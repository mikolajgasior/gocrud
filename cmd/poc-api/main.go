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
	"strconv"
	"syscall"
	"time"

	"codeberg.org/mikolajgasior/gocrud"
	crudapi "codeberg.org/mikolajgasior/gocrud/pkg/http/api"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
	_ "modernc.org/sqlite"
)

// ── Structs ───────────────────────────────────────────────────────────────────

type User struct {
	ID         uint64
	Name       string `crud:"req len:1,100"`
	Email      string `crud:"req email uniq"`
	CreatedAt  int64
	CreatedBy  uint64
	ModifiedAt int64
	ModifiedBy uint64
}

// User_Draft is used for the create operation — callers only provide Name and
// Email. The audit fields are included so the service layer populates them
// automatically; callers never need to send them.
type User_Draft struct {
	ID         uint64
	Name       string `crud:"req len:1,100"`
	Email      string `crud:"req email uniq"`
	CreatedAt  int64
	CreatedBy  uint64
	ModifiedAt int64
	ModifiedBy uint64
}

type Note struct {
	ID         uint64
	Title      string `crud:"req len:1,200"`
	Content    string
	UserID     uint64 `crud:"req"`
	CreatedAt  int64
	CreatedBy  uint64
	ModifiedAt int64
	ModifiedBy uint64
}

// Note_Draft is used for the create operation — callers only provide Title,
// Content, and UserID. Audit fields are included so the service layer
// populates them automatically.
type Note_Draft struct {
	ID         uint64
	Title      string `crud:"req len:1,200"`
	Content    string
	UserID     uint64 `crud:"req"`
	CreatedAt  int64
	CreatedBy  uint64
	ModifiedAt int64
	ModifiedBy uint64
}

// ── Main ──────────────────────────────────────────────────────────────────────

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
		"users": func() interface{} { return &User{} },
		"notes": func() interface{} { return &Note{} },
	}

	crud := gocrud.New(db, gocrud.Options{Dialect: gocrud.DialectSQLite})
	svc := svccrud.New(paths, db, gocrud.DialectSQLite)

	if err := crud.CreateTable(context.Background(), &User{}); err != nil {
		slog.Error("failed to create users table", "error", err)
		os.Exit(1)
	}
	if err := crud.CreateTable(context.Background(), &Note{}); err != nil {
		slog.Error("failed to create notes table", "error", err)
		os.Exit(1)
	}

	// ── HTTP API ──────────────────────────────────────────────────────────────

	// The caller identifies themselves via the X-User-ID header.
	// This is intentionally simple — a real application would use a proper
	// authentication mechanism (JWT, session, etc.).
	userIDFromRequest := func(r *http.Request) uint64 {
		id, _ := strconv.ParseUint(r.Header.Get("X-User-ID"), 10, 64)
		return id
	}

	handler := crudapi.New(svc, crudapi.Options{
		UserIDFunc: userIDFromRequest,
		Paths: map[string]crudapi.PathOptions{
			"users": {
				CreateConstructor: func() interface{} { return &User_Draft{} },
			},
			"notes": {
				CreateConstructor: func() interface{} { return &Note_Draft{} },
				AllowedFilters:    []string{"UserID"},
			},
		},
	})

	mux := http.NewServeMux()
	subMux := http.NewServeMux()
	subMux.HandleFunc("/", handler.Serve)
	mux.Handle("/api/", http.StripPrefix("/api", subMux))

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
		slog.Info("users API",
			"list",   "GET    /api/users/",
			"read",   "GET    /api/users/{id}",
			"create", "PUT    /api/users/",
			"update", "PUT    /api/users/{id}",
			"delete", "DELETE /api/users/{id}",
		)
		slog.Info("notes API",
			"list",         "GET    /api/notes/",
			"list_by_user", "GET    /api/notes/?filter_val_UserID={id}&filter_op_UserID=eq",
			"read",         "GET    /api/notes/{id}",
			"create",       "PUT    /api/notes/  (X-User-ID header sets CreatedBy/ModifiedBy)",
			"update",       "PUT    /api/notes/{id}",
			"delete",       "DELETE /api/notes/{id}",
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
