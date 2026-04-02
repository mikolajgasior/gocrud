package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	"codeberg.org/mikolajgasior/gocrud/pkg/poc/dbconn"
)

func portFromEnvOrExit() string {
	port := os.Getenv("PORT")
	_, err := strconv.Atoi(port)
	if err != nil {
		slog.Error(LogInvalidPort, "port", port)
		os.Exit(ExitInvalidPort)
	}
	slog.Info("port number from env", "port", port)

	return port
}

func createTablesEnabled() bool {
	return os.Getenv("CREATE_TABLES") == "true"
}

func dbConnOrExit() *sql.DB {
	dbConn, err := dbconn.NewFromEnv()
	if err != nil {
		slog.Error(LogCannotOpenDB, logger.AttrError(err))
		os.Exit(ExitCannotOpenDB)
	}

	return dbConn
}

func timeoutMiddleware(t time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), t)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
