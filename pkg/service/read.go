package service

import (
	"context"
	"fmt"
	"log/slog"

	"codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

// Read loads a single record by id. When constructor is nil the path's
// registered constructor is used; pass a non-nil value to load into a
// different struct type (e.g. a read-specific projection).
func (c *CRUD) Read(ctx context.Context, path string, id uint64, constructor func() interface{}) (interface{}, error) {
	logAttrService := logger.AttrService(c, "Read")

	if constructor == nil {
		var ok bool
		constructor, ok = c.paths[path]
		if !ok {
			slog.Error("path not found", logAttrService)
			return nil, InvalidPathError
		}
	}

	obj := constructor()
	err := c.crud.Load(ctx, obj, fmt.Sprintf("%d", id), gocrud.LoadOptions{})
	if err != nil {
		return nil, err
	}

	objID := gocrud.ObjIDValue(obj)
	if objID == 0 {
		slog.Error("error not found", logAttrService, slog.Uint64("id", id))
		return nil, NotFoundError
	}

	return obj, nil
}
