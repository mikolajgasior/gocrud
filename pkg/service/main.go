package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"codeberg.org/mikolajgasior/gocrud"
	sqlfilters "codeberg.org/mikolajgasior/gocrud/pkg/filters"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	validator "github.com/mikolajgasior/struct-validator"
)

var (
	InvalidPathError = errors.New("invalid path")
)

type CRUD struct {
	paths map[string]func() interface{}
	crud  *gocrud.CRUD
}

func New(paths map[string]func() interface{}, dbConn *sql.DB, dialect string) *CRUD {
	return &CRUD{
		paths: paths,
		crud:  gocrud.New(dbConn, gocrud.Options{Dialect: dialect}),
	}
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

func (c *CRUD) ID(obj interface{}) uint64 {
	return gocrud.ObjIDValue(obj)
}

// Paths returns the registered path keys.
func (c *CRUD) Paths() []string {
	keys := make([]string, 0, len(c.paths))
	for k := range c.paths {
		keys = append(keys, k)
	}
	return keys
}

// buildFilters converts the string-keyed filter maps coming from HTTP layers
// into a sqlfilters.Filters value ready for the CRUD layer.
func buildFilters(filterVals, filterOps map[string]string, logAttr slog.Attr) (*sqlfilters.Filters, error) {
	filters := sqlfilters.Filters{}
	violations := map[string]uint64{}

	for name, value := range filterVals {
		op := sqlfilters.OpEqual
		if filterOp, ok := filterOps[name]; ok {
			var err error
			op, err = Op(filterOp)
			if err != nil {
				slog.Error("invalid filter op", logAttr, slog.String("op", filterOp))
				violations[name] = validator.FailType
			}
		}
		filters[name] = sqlfilters.OpVal{Op: op, Val: value}
	}

	if len(violations) > 0 {
		slog.Error("invalid filter ops", logAttr)
		return nil, &FilterValidationError{
			Err:        errors.New("invalid filter ops"),
			Violations: violations,
		}
	}

	return &filters, nil
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
