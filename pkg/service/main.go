package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	structcrud "codeberg.org/mikolajgasior/gocrud"
	sqlfilters "codeberg.org/mikolajgasior/gocrud/pkg/filters"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

var (
	InvalidPathError = errors.New("invalid path")
)

type CRUD struct {
	paths map[string]func() interface{}
	crud  *structcrud.CRUD
}

func New(paths map[string]func() interface{}, dbConn *sql.DB) *CRUD {
	crud := &CRUD{
		paths: paths,
		crud:  structcrud.New(dbConn, structcrud.Options{}),
	}

	return crud
}

func (c *CRUD) CreateTables(ctx context.Context) error {
	logAttrService := logger.AttrService(c, "CreateTables")

	var err error
	for _, constructor := range c.paths {
		err = c.crud.CreateTable(ctx, constructor())
		if err != nil {
			slog.Error("error creating table", logAttrService, logger.AttrError(err))
			return err
		}
	}
	return nil
}

func (c *CRUD) New(path string) interface{} {
	constructor, ok := c.paths[path]
	if !ok {
		return nil
	}
	return constructor()
}

func (c *CRUD) ID(obj interface{}) int64 {
	return structcrud.ObjIDValue(obj)
}

func Op(op string) (int, error) {
	switch op {
	case "eq":
		return sqlfilters.OpEqual, nil
	case "ne":
		return sqlfilters.OpNotEqual, nil
	case "lt":
		return sqlfilters.OpLower, nil
	case "le":
		return sqlfilters.OpLowerOrEqual, nil
	case "gt":
		return sqlfilters.OpGreater, nil
	case "ge":
		return sqlfilters.OpGreaterOrEqual, nil
	case "like":
		return sqlfilters.OpLike, nil
	case "match":
		return sqlfilters.OpMatch, nil
	case "bit":
		return sqlfilters.OpBit, nil
	default:
		return 0, InvalidFilterOpError
	}
}
