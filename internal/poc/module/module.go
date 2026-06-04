package module

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/layout"
)

const (
	LogCannotCreateDBTable = "failed to create table"
	LogLayoutIsNil         = "layout is nil"
)

var (
	LayoutIsNilError = errors.New(LogLayoutIsNil)
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
	Sitemap() *layout.Sitemap
}

type InitInput struct {
	DBConn       *sql.DB
	CreateTables bool
	Dialect      string
}
