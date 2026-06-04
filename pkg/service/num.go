package service

import (
	"context"
	"errors"
	"log/slog"

	"codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

func (c *CRUD) Num(ctx context.Context, path string, filterVals, filterOps map[string]string) (uint64, error) {
	logAttrService := logger.AttrService(c, "Num")

	constructor, ok := c.paths[path]
	if !ok {
		slog.Error("path not found", logAttrService)
		return 0, InvalidPathError
	}

	filters, err := buildFilters(filterVals, filterOps, logAttrService)
	if err != nil {
		return 0, err
	}

	numObjs, err := c.crud.GetCount(ctx, constructor(), crud.GetCountOptions{
		Filters:                  filters,
		ConvertFiltersFromString: true,
	})

	if err != nil {
		var crudErr *crud.CRUDError
		if errors.As(err, &crudErr) {
			var validationErr *crud.ValidationError
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
