package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mikolajgasior/gocrud"
	"github.com/mikolajgasior/gocrud/pkg/logger"
)

// Read loads a single record by id. When constructor is nil the registry's
// constructor for key is used; pass a non-nil value to load into a
// different struct type (e.g. a read-specific projection).
func (c *CRUD) Read(ctx context.Context, key string, id uint64, constructor func() interface{}) (interface{}, error) {
	logAttrService := logger.AttrService(c, "Read")

	if constructor == nil {
		var ok bool
		constructor, ok = c.registry[key]
		if !ok {
			slog.Error("key not found", logAttrService)
			return nil, InvalidKeyError
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
