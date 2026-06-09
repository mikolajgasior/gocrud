package service

import (
	"context"
	"errors"
	"log/slog"

	"codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

// List returns a paginated, filtered list of records. When constructor is nil
// the registry's constructor for key is used; pass a non-nil value to scan rows
// into a different struct type (e.g. a list-specific projection).
func (c *CRUD) List(ctx context.Context, key string, limit, offset int, order, orderDirection string, filterVals, filterOps map[string]string, rowFunc func(interface{}) interface{}, constructor func() interface{}) ([]interface{}, error) {
	logAttrService := logger.AttrService(c, "List")

	if constructor == nil {
		var ok bool
		constructor, ok = c.registry[key]
		if !ok {
			slog.Error("key not found", logAttrService)
			return nil, InvalidKeyError
		}
	}

	filters, err := buildFilters(filterVals, filterOps, logAttrService)
	if err != nil {
		return nil, err
	}

	getOrder := []string{}
	if order != "" {
		getOrder = append(getOrder, order, orderDirection)
	}

	objs, err := c.crud.Get(ctx, constructor, gocrud.GetOptions{
		Limit:                    limit,
		Offset:                   offset,
		Order:                    getOrder,
		Filters:                  filters,
		ConvertFiltersFromString: true,
		RowObjTransformFunc:      rowFunc,
	})

	if err != nil {
		var crudErr *gocrud.CRUDError
		if errors.As(err, &crudErr) {
			var validationErr *gocrud.ValidationError
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
