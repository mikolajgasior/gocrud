package service

import (
	"context"
	"log/slog"

	"codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

func (c *CRUD) Delete(ctx context.Context, path string, id uint64) error {
	logAttrService := logger.AttrService(c, "Delete")

	obj, err := c.Read(ctx, path, id, nil)
	if err != nil {
		return err
	}

	err = c.crud.Delete(ctx, obj, gocrud.DeleteOptions{})
	if err != nil {
		slog.Error("error deleting", logAttrService, slog.Uint64("id", id))
		return err
	}

	return nil
}
