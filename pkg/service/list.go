package service

import (
	"context"
	"errors"
	"log/slog"

	structcrud "miko.gs/gocrud"
	sqlfilters "miko.gs/gocrud/pkg/filters"
	"miko.gs/gocrud/pkg/logger"
	validator "miko.gs/struct-validator"
)

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
