package dbconn

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type DBConnError struct {
	Op  string
	Tag string
	Err error
}

func (e *DBConnError) Error() string {
	return e.Op + ": " + e.Err.Error()
}

func NewFromEnv() (*sql.DB, error) {
	vars, err := connEnvVars()
	if err != nil {
		return nil, &DBConnError{
			Op:  "get connection env vars",
			Err: err,
		}
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		vars["DB_USER"], vars["DB_PASS"], vars["DB_HOST"], vars["DB_PORT"], vars["DB_NAME"], "disable",
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, &DBConnError{
			Op:  "open db connection",
			Err: err,
		}
	}
	if err = dbConn.Ping(); err != nil {
		return nil, &DBConnError{
			Op:  "ping db connection",
			Err: err,
		}
	}

	return dbConn, nil
}

func connEnvVars() (map[string]string, error) {
	dbEnvs := make(map[string]string, 5)
	for _, key := range []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME",
	} {
		val := os.Getenv(key)
		if val == "" {
			return nil, fmt.Errorf("empty %s", key)
		}
		dbEnvs[key] = val
	}

	_, err := strconv.Atoi(dbEnvs["DB_PORT"])
	if err != nil {
		return nil, fmt.Errorf("DB_PORT must be a numeric value: %v", err)
	}

	return dbEnvs, nil
}
