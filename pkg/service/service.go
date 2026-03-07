package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	sqlbuilder "miko.gs/pgsql-builder"
	sqlfilters "miko.gs/pgsql-builder/pkg/filters"
	structcrud "miko.gs/struct-crud"
	"miko.gs/struct-crud/pkg/logger"
	validator "miko.gs/struct-validator"
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

func (c *CRUD) Read(ctx context.Context, path string, id int64) (interface{}, error) {
	logAttrService := logger.AttrService(c, "Read")

	constructor, ok := c.paths[path]
	if !ok {
		slog.Error("path not found", logAttrService)
		return nil, InvalidPathError
	}

	obj := constructor()
	err := c.crud.Load(ctx, obj, fmt.Sprintf("%d", id), structcrud.LoadOptions{})
	if err != nil {
		return nil, err
	}

	objID := structcrud.ObjIDValue(obj)
	if objID == 0 {
		slog.Error("error not found", logAttrService, slog.Int64("id", objID))
		return nil, NotFoundError
	}

	return obj, nil
}

func (c *CRUD) Save(ctx context.Context, obj interface{}, now, userID int64) error {
	logAttrService := logger.AttrService(c, "Save")

	err := c.crud.Save(ctx, obj, structcrud.SaveOptions{
		ModifiedAt: now,
		ModifiedBy: userID,
	})
	if err != nil {
		var crudErr *structcrud.CRUDError
		if errors.As(err, &crudErr) {
			var validationErr *structcrud.ValidationError
			if errors.As(crudErr.Err, &validationErr) {
				slog.Error("error with validation", logAttrService, logger.AttrError(err))
				return &ModelValidationError{
					Err: validationErr,
				}
			}
		}

		slog.Error("error saving", logAttrService, logger.AttrError(err))
		return err
	}

	return nil
}

func (c *CRUD) SaveFromForm(ctx context.Context, obj interface{}, values url.Values, namePrefix string, now, userID int64) error {
	logAttrService := logger.AttrService(c, "SaveFromForm")

	fieldViolations := map[string]int{}
	for key, value := range values {
		name := strings.Replace(key, namePrefix, "", 1)

		ok, _ := sqlbuilder.SetStructFieldValueFromString(obj, name, value[0])
		if !ok {
			fieldViolations[name] = validator.FailType
		}
	}

	if len(fieldViolations) > 0 {
		slog.Error("error converting form value to obj field value", logAttrService)
		return &ModelValidationError{
			Err:        errors.New("invalid value type"),
			Violations: fieldViolations,
		}
	}

	err := c.crud.Save(ctx, obj, structcrud.SaveOptions{
		ModifiedAt: now,
		ModifiedBy: userID,
	})
	if err != nil {
		var crudErr *structcrud.CRUDError
		if errors.As(err, &crudErr) {
			var validationErr *structcrud.ValidationError
			if errors.As(crudErr.Err, &validationErr) {
				slog.Error("error with validation", logAttrService, logger.AttrError(err))
				return &ModelValidationError{
					Err:        validationErr,
					Violations: validationErr.Violations,
				}
			}
		}

		slog.Error("error saving", logAttrService, logger.AttrError(err))
		return err
	}

	return nil
}

func (c *CRUD) Delete(ctx context.Context, path string, id int64) error {
	logAttrService := logger.AttrService(c, "Delete")

	obj, err := c.Read(ctx, path, id)
	if err != nil {
		return err
	}

	err = c.crud.Delete(ctx, obj, structcrud.DeleteOptions{})
	if err != nil {
		slog.Error("error deleting", logAttrService, slog.Int64("id", id))
		return err
	}

	return nil
}

func (c *CRUD) List(ctx context.Context, path string, limit, offset int, order, orderDirection string, filterVals, filterOps map[string]string, rowFunc func(interface{}) interface{}) ([]interface{}, error) {
	logAttrService := logger.AttrService(c, "List")

	filterViolations := map[string]int{}
	filters := sqlfilters.Filters{}
	for name, value := range filterVals {
		op := sqlfilters.OpEqual
		filterOp, ok := filterOps[name]
		if ok {
			var err error
			op, err = Op(filterOp)
			if err != nil {
				slog.Error("invalid filter op", logAttrService, slog.String("op", filterOp))
				filterViolations[name] = validator.FailType
			}
		}

		filters[name] = sqlfilters.OpVal{
			Op:  op,
			Val: value,
		}
	}

	if len(filterViolations) > 0 {
		slog.Error("invalid filter ops", logAttrService)
		return nil, &FilterValidationError{
			Err:        errors.New("invalid filter ops"),
			Violations: filterViolations,
		}
	}

	getOrder := []string{}
	if order != "" {
		getOrder = append(getOrder, order, orderDirection)
	}

	objs, err := c.crud.Get(ctx, c.paths[path], structcrud.GetOptions{
		Limit:                    limit,
		Offset:                   offset,
		Order:                    getOrder,
		Filters:                  &filters,
		ConvertFiltersFromString: true,
		RowObjTransformFunc:      rowFunc,
	})

	if err != nil {
		var crudErr *structcrud.CRUDError
		if errors.As(err, &crudErr) {
			var validationErr *structcrud.ValidationError
			if errors.As(crudErr.Err, &validationErr) {
				slog.Error("error with validation", logAttrService, logger.AttrError(validationErr))
				return nil, &FilterValidationError{
					Err:        validationErr,
					Violations: validationErr.Violations,
				}
			}
		}

		slog.Error("error getting", logAttrService, logger.AttrError(err))
		return nil, err
	}

	return objs, nil
}

func (c *CRUD) Num(ctx context.Context, path string, limit, offset int, order, orderDirection string, filterVals, filterOps map[string]string) (int64, error) {
	logAttrService := logger.AttrService(c, "Num")

	filterViolations := map[string]int{}
	filters := sqlfilters.Filters{}
	for name, value := range filterVals {
		op := sqlfilters.OpEqual
		filterOp, ok := filterOps[name]
		if ok {
			var err error
			op, err = Op(filterOp)
			if err != nil {
				slog.Error("invalid filter op", logAttrService, slog.String("op", filterOp))
				filterViolations[name] = validator.FailType
			}
		}

		filters[name] = sqlfilters.OpVal{
			Op:  op,
			Val: value,
		}
	}

	if len(filterViolations) > 0 {
		slog.Error("invalid filter ops", logAttrService)
		return 0, &FilterValidationError{
			Err:        errors.New("invalid filter ops"),
			Violations: filterViolations,
		}
	}

	getOrder := []string{}
	if order != "" {
		getOrder = append(getOrder, order, orderDirection)
	}

	obj := c.paths[path]()

	numObjs, err := c.crud.GetCount(ctx, obj, structcrud.GetCountOptions{
		Filters:                  &filters,
		ConvertFiltersFromString: true,
	})

	if err != nil {
		var crudErr *structcrud.CRUDError
		if errors.As(err, &crudErr) {
			var validationErr *structcrud.ValidationError
			if errors.As(crudErr.Err, &validationErr) {
				slog.Error("error with validation", logAttrService, logger.AttrError(validationErr))
				return 0, &FilterValidationError{
					Err:        validationErr,
					Violations: validationErr.Violations,
				}
			}
		}

		slog.Error("error getting", logAttrService, logger.AttrError(err))
		return 0, err
	}

	return numObjs, nil
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
