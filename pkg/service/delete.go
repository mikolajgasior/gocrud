package service

import (
	"context"
	"log/slog"

	structcrud "miko.gs/gocrud"
	"miko.gs/gocrud/pkg/logger"
)

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
