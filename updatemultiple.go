package crud

import (
	"context"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
	validator "github.com/keenbytes/struct-validator"
)

type UpdateMultipleOptions struct {
	Filters                 *sqlbuilder.Filters
	ConvertValuesFromString bool
}

func (c *CRUD) UpdateMultiple(ctx context.Context, obj interface{}, fieldsToUpdate map[string]interface{}, options UpdateMultipleOptions) error {
	builder, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}
	if options.ConvertValuesFromString {
		for name, value := range fieldsToUpdate {
			valueAsString, ok := value.(string)
			if !ok {
				continue
			}

			ok, valueAsFieldType := sqlbuilder.StructFieldValueFromString(obj, name, valueAsString)
			if !ok {
				return getStructFieldValueFromStringCRUDError()
			}
			fieldsToUpdate[name] = valueAsFieldType
		}
	}
	restrictFields := make(map[string]bool, len(fieldsToUpdate))
	for name := range fieldsToUpdate {
		restrictFields[name] = true
	}

	ok, violations, err := validator.Validate(obj, &validator.ValidateOptions{
		RestrictFields:  restrictFields,
		OverwriteValues: fieldsToUpdate,
		TagName:         c.tagName,
	})

	if err != nil {
		return getValidateObjCRUDError(err)
	}

	if !ok {
		return getObjInvalidCRUDError(violations)
	}

	err = ValidateFilters(obj, options.Filters, c.tagName)
	if err != nil {
		return err
	}

	var query string
	// if the object has a Update method, use it.
	if updateerImpl, ok := obj.(updateer); ok {
		query, err = updateerImpl.Update(fieldsToUpdate, options.Filters)
		if err != nil {
			return getObjFuncCRUDError("update", err)
		}
	} else {
		query, err = builder.Update(fieldsToUpdate, options.Filters)
		if err != nil {
			return getBuilderFuncCRUDError("update", err)
		}
	}

	_, err = c.db.ExecContext(ctx, query, append(sqlbuilder.MapInterfaces(fieldsToUpdate), sqlbuilder.FiltersInterfaces(options.Filters)...)...)
	if err != nil {
		return getDBFuncCRUDError("exec", err)
	}

	return nil
}
