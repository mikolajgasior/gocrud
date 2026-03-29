package service

import (
	"context"
	"errors"
	"log/slog"

	structcrud "codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

func (c *CRUD) Save(ctx context.Context, obj interface{}, now int64, userID uint64) error {
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
