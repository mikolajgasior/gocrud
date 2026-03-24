package service

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"strings"

	"miko.gs/gocrud/pkg/logger"
	sqlbuilder "miko.gs/gocrud/pkg/sql/builder"
	validator "miko.gs/struct-validator"
)

func (c *CRUD) SaveFromForm(ctx context.Context, obj interface{}, values url.Values, namePrefix string, now, userID int64) error {
	logAttrService := logger.AttrService(c, "SaveFromForm")

	fieldViolations := map[string]int{}
	for key, value := range values {
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
