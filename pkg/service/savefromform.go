package service

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"strings"

	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
	sqlbuilder "codeberg.org/mikolajgasior/gocrud/pkg/sql/builder"
	validator "github.com/mikolajgasior/struct-validator"
)

func (c *CRUD) SaveFromForm(ctx context.Context, obj interface{}, values url.Values, namePrefix string, now int64, userID uint64) error {
	logAttrService := logger.AttrService(c, "SaveFromForm")

	fieldViolations := map[string]uint64{}
	for key, value := range values {
		if len(value) == 0 {
			continue
		}
		name := strings.Replace(key, namePrefix, "", 1)

		ok, _ := sqlbuilder.SetStructFieldValueFromString(obj, name, value[0])
		if !ok {
			fieldViolations[name] = validator.FailType
		}
	}

	if len(fieldViolations) > 0 {
		slog.Error("error converting form value to obj field value", logAttrService)
		return &ModelValidationError{
			Err:        errors.New("invalid value type"),
			Violations: fieldViolations,
		}
	}

	return c.Save(ctx, obj, now, userID)
}
