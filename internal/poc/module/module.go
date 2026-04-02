package module

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
)

const (
	LogCannotCreateDBTable = "failed to create table"
)

type InitError struct {
	ExitCode int
}

func (e *InitError) Error() string {
	return "exit code: " + strconv.Itoa(e.ExitCode)
}

type CreateTableError struct {
	Err error
}

func (e *CreateTableError) Error() string {
	return "failed to create table: " + e.Err.Error()
}

type Module interface {
	Init(context.Context, InitInput) error
	AddHandler(*http.ServeMux)
}

type InitInput struct {
	DBConn       *sql.DB
	CreateTables bool
}
