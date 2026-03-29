package service

import (
	"context"
	"errors"
	"log/slog"

	structcrud "codeberg.org/mikolajgasior/gocrud"
	sqlfilters "codeberg.org/mikolajgasior/gocrud/pkg/filters"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	validator "github.com/mikolajgasior/struct-validator"
)

func (c *CRUD) Num(ctx context.Context, path string, filterVals, filterOps map[string]string) (uint64, error) {
	logAttrService := logger.AttrService(c, "Num")

	filterViolations := map[string]uint64{}
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
