package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codeberg.org/mikolajgasior/gocrud"
	crudapi "codeberg.org/mikolajgasior/gocrud/pkg/http/api"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
	_ "modernc.org/sqlite"
)

// Note is the primary struct — maps to the "note" table.
type Note struct {
	ID         uint64
	Title      string `crud:"req len:1,200"`
	Content    string
	CreatedAt  int64
	CreatedBy  uint64
	ModifiedAt int64
	ModifiedBy uint64
}

// Note_Draft is used for the create operation.
// The "_" suffix is stripped when deriving the table name so it still
// maps to the "note" table, but only Title and Content are required —
// the audit fields are populated by the service layer.
type Note_Draft struct {
	ID      uint64
	Title   string `crud:"req len:1,200"`
	Content string
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "poc-api.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ── Database ──────────────────────────────────────────────────────────────

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// ── CRUD + Service ────────────────────────────────────────────────────────

	paths := map[string]func() interface{}{
		"notes": func() interface{} { return &Note{} },
	}

	crud := gocrud.New(db, gocrud.Options{Dialect: gocrud.DialectSQLite})
	svc := svccrud.New(paths, db, gocrud.DialectSQLite)

	// Create table if it doesn't exist.
	if err := crud.CreateTable(context.Background(), &Note{}); err != nil {
		slog.Error("failed to create notes table", "error", err)
		os.Exit(1)
	}

	// ── HTTP API ──────────────────────────────────────────────────────────────

	handler := crudapi.New(svc, crudapi.Options{
		Paths: map[string]crudapi.PathOptions{
			"notes": {
				// Use the leaner Note_Draft struct for creates so callers
				// don't have to send audit fields.
				CreateConstructor: func() interface{} { return &Note_Draft{} },
			},
		},
	})

	mux := http.NewServeMux()
	subMux := http.NewServeMux()
	subMux.HandleFunc("/", handler.Serve)
	mux.Handle("/api/", http.StripPrefix("/api", subMux))

	// Simple health-check endpoint.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	// ── Server lifecycle ──────────────────────────────────────────────────────

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("starting server", "port", port, "db", dbPath)
		slog.Info("endpoints",
			"list",   "GET  /api/notes/",
			"read",   "GET  /api/notes/{id}",
			"create", "PUT  /api/notes/",
			"update", "PUT  /api/notes/{id}",
			"delete", "DELETE /api/notes/{id}",
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("shutting down")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
