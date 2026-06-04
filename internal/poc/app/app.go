package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/module"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

const (
	requestTimeoutSeconds = 5
)

type App struct {
	Modules map[string]module.Module
}

func (a *App) Run(ctx context.Context) {
	logger.SetLogger(os.Getenv("LOG_LEVEL"))

	// main http server context that handles signals
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,  // Ctrl‑C
		syscall.SIGTERM, // termination request
	)
	defer stop()

	port := portFromEnvOrExit()
	createTables := createTablesEnabled()

	dbConn := dbConnOrExit()
	defer func() {
		err := dbConn.Close()
		if err != nil {
			slog.Error(LogCannotCloseDB, logger.AttrError(err))
			os.Exit(ExitCannotCloseDB)
		}
	}()

	moduleNames := make([]string, 0, len(a.Modules))
	for name := range a.Modules {
		moduleNames = append(moduleNames, name)
	}
	sort.Strings(moduleNames)

	for _, name := range moduleNames {
		mod := a.Modules[name]
		err := mod.Init(ctx, module.InitInput{
			DBConn:       dbConn,
			CreateTables: createTables,
			Dialect:      "postgres",
		})
		if err != nil {
			var errModInit *module.InitError
			if errors.As(err, &errModInit) {
				os.Exit(errModInit.ExitCode)
			}

			var errModCreateTable *module.CreateTableError
			if errors.As(err, &errModCreateTable) {
				os.Exit(ExitCreateTableError)
			}

			os.Exit(ExitInitError)
		}
	}

	serveMux := http.NewServeMux()
	for _, name := range moduleNames {
		mod := a.Modules[name]
		mod.AddHandler(serveMux)
	}

	// apply a timeout to every request.
	wrapped := timeoutMiddleware(requestTimeoutSeconds * time.Second)(serveMux)
	slog.Info("applied timeout for requests", "timeout", fmt.Sprintf("%ds", requestTimeoutSeconds))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: wrapped,
	}

	go func() {
		slog.Info("starting http server", "port", port)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			slog.Error(LogHTTPListen, logger.AttrError(err))
			os.Exit(ExitHTTPListen)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("shutting down http server")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed - forcing close", logger.AttrError(err))
		if cerr := srv.Close(); cerr != nil {
			slog.Error("forced close also failed", logger.AttrError(cerr))
		}
	}

	slog.Info("http server stopped")
}
