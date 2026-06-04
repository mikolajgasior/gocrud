package service

import (
	"context"
	"errors"
	"log/slog"

	structcrud "codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

func (c *CRUD) List(ctx context.Context, path string, limit, offset int, order, orderDirection string, filterVals, filterOps map[string]string, rowFunc func(interface{}) interface{}) ([]interface{}, error) {
	logAttrService := logger.AttrService(c, "List")

	constructor, ok := c.paths[path]
	if !ok {
		slog.Error("path not found", logAttrService)
		return nil, InvalidPathError
	}

	filters, err := buildFilters(filterVals, filterOps, logAttrService)
	if err != nil {
		return nil, err
	}

	getOrder := []string{}
	if order != "" {
		getOrder = append(getOrder, order, orderDirection)
	}

	objs, err := c.crud.Get(ctx, constructor, structcrud.GetOptions{
		Limit:                    limit,
		Offset:                   offset,
		Order:                    getOrder,
		Filters:                  filters,
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
