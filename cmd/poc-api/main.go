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

type User struct {
	ID         uint64
	Name       string `crud:"req len:1,100"`
	Email      string `crud:"req email uniq"`
	CreatedAt  int64
	CreatedBy  uint64
	ModifiedAt int64
	ModifiedBy uint64
}

// User_Draft is the create payload; the underscore tells gocrud to use the "user" table.
type User_Draft struct {
	ID    uint64
	Name  string `crud:"req len:1,100"`
	Email string `crud:"req email uniq"`
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

// Note_Draft is the create payload; the underscore tells gocrud to use the "note" table.
type Note_Draft struct {
	ID      uint64
	Title   string `crud:"req len:1,200"`
	Content string
	UserID  uint64 `crud:"req"`
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

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	svc := svccrud.New(map[string]func() interface{}{
		"users": func() interface{} { return &User{} },
		"notes": func() interface{} { return &Note{} },
	}, db, gocrud.DialectSQLite)

	if err := svc.CreateTables(context.Background()); err != nil {
		slog.Error("failed to create tables", "error", err)
		os.Exit(1)
	}

	// X-User-ID is intentionally simple — a real app would use JWT or sessions.
	userIDFromRequest := func(r *http.Request) uint64 {
		id, _ := strconv.ParseUint(r.Header.Get("X-User-ID"), 10, 64)
		return id
	}

	handler := crudapi.New(svc, crudapi.Options{
		UserIDFunc: userIDFromRequest,
		Routes: map[string]crudapi.Route{
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
	mux.Handle("/api/", http.StripPrefix("/api", http.HandlerFunc(handler.Serve)))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("starting server", "port", port, "db", dbPath)
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
