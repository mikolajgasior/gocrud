package crud

import (
	"errors"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
	validator "github.com/keenbytes/struct-validator"
)

type UpdateMultipleOptions struct {
	Filters                 *sqlbuilder.Filters
	ConvertValuesFromString bool
}

func (c *CRUD) UpdateMultiple(obj interface{}, fieldsToUpdate map[string]interface{}, options UpdateMultipleOptions) error {
	builder, err := c.builder(obj)
	if err != nil {
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}
	if options.ConvertValuesFromString {
		for name, value := range fieldsToUpdate {
			valueAsString, ok := value.(string)
			if !ok {
				continue
			}

			ok, valueAsFieldType := sqlbuilder.StructFieldValueFromString(obj, name, valueAsString)
			if !ok {
				return ErrCRUD{
					Op:  "sqlbuilder.StructFieldValueFromString",
					Err: errors.New("error converting string to field type"),
				}
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
		return ErrCRUD{
			Op:  "Validate",
			Err: err,
		}
	}

	if !ok {
		return ErrCRUD{
			Op: "Validate",
			Err: &ErrValidation{
				Violations: violations,
			},
		}
	}

	err = ValidateFilters(obj, options.Filters, c.tagName)
	if err != nil {
		return err
	}

	query, err := builder.Update(fieldsToUpdate, options.Filters)
	if err != nil {
		return ErrCRUD{
			Op:  "builder.Update",
			Err: err,
		}
	}

	_, err = c.db.Exec(query, append(sqlbuilder.MapInterfaces(fieldsToUpdate), sqlbuilder.FiltersInterfaces(options.Filters)...)...)
	if err != nil {
		return ErrCRUD{
			Op:  "o.db.Exec",
			Err: err,
		}
	}

	return nil
}
