package service

import (
	"context"
	"fmt"
	"log/slog"

	"codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

func (c *CRUD) Read(ctx context.Context, path string, id uint64) (interface{}, error) {
	logAttrService := logger.AttrService(c, "Read")

	constructor, ok := c.paths[path]
	if !ok {
		slog.Error("path not found", logAttrService)
		return nil, InvalidPathError
	}

	obj := constructor()
	err := c.crud.Load(ctx, obj, fmt.Sprintf("%d", id), gocrud.LoadOptions{})
	if err != nil {
		return nil, err
	}

	objID := gocrud.ObjIDValue(obj)
	if objID == 0 {
		slog.Error("error not found", logAttrService, slog.Uint64("id", objID))
		return nil, NotFoundError
	}

	return obj, nil
}
