package gocrud

import (
	sqlfilters "codeberg.org/mikolajgasior/gocrud/pkg/filters"
	validator "github.com/mikolajgasior/struct-validator"
)

func Validate(obj interface{}, restrictFields map[string]bool, tagName string) (bool, map[string]uint64, error) {
	ok, violations, err := validator.Validate(obj, &validator.ValidateOptions{
		RestrictFields: restrictFields,
		TagName:        tagName,
	})

	return ok, violations, err
}

func ValidateFilters(obj interface{}, filters *sqlfilters.Filters, tagName string) error {
	if filters == nil || len(*filters) == 0 {
		return nil
	}

	err := sqlfilters.SetObjFields(obj, filters)
	if err != nil {
		return getBuilderFuncCRUDError("set obj fields", err)
	}

	fieldsToValidate := map[string]bool{}
	for filter := range *filters {
		fieldsToValidate[filter] = true
	}

	ok, violations, err := Validate(obj, fieldsToValidate, tagName)
	if err != nil {
		return getValidateObjCRUDError(err)
	}

	if !ok {
		return getObjInvalidCRUDError(violations)
	}

	return nil
}
