package crud

import (
	"context"

	sqlbuilder "miko.gs/pgsql-builder"
	sqlfilters "miko.gs/pgsql-builder/pkg/filters"
	validator "miko.gs/struct-validator"
)

type GetCountOptions struct {
	Filters                  *sqlfilters.Filters
	ConvertFiltersFromString bool
}

// GetCount runs a 'SELECT COUNT(*)' query on the database with specified filters, order, limit and offset and returns count of rows
func (c *CRUD) GetCount(ctx context.Context, obj interface{}, options GetCountOptions) (int64, error) {
	builder, err := c.builder(obj)
	if err != nil {
		return 0, getBuilderObjectCRUDError(err)
	}

	// Filter values can be passed as string. We do not want any use of reflect outside of CRUD.
	if options.ConvertFiltersFromString {
		newFilters := &sqlfilters.Filters{}

		for filterName, filterOpVal := range *(options.Filters) {
			// ignore the raw filters entirely, that's too complicated
			if filterName == sqlfilters.Raw {
				continue
			}

			valueAsString, ok := filterOpVal.Val.(string)
			if !ok {
				return 0, getObjInvalidCRUDError(map[string]int{
					filterName: validator.FailType,
				})
			}

			ok, valueAsFieldType := sqlbuilder.StructFieldValueFromString(obj, filterName, valueAsString)
			if !ok {
				return 0, getObjInvalidCRUDError(map[string]int{
					filterName: validator.FailType,
				})
			}

			newFilters.Add(filterName, sqlfilters.OpVal{
				Op:  filterOpVal.Op,
				Val: valueAsFieldType,
			})
		}

		options.Filters = newFilters
	}

	err = ValidateFilters(obj, options.Filters, c.tagName)
	if err != nil {
		return 0, err
	}

	var query string
	// if the object has a SelectCount method, use it.
	if selectCounterImpl, ok := obj.(selectCountQueryBuilder); ok {
		query, err = selectCounterImpl.SelectCountQuery(options.Filters)
		if err != nil {
			return 0, getObjFuncCRUDError("select count query", err)
		}
	} else {
		query, err = builder.SelectCount(options.Filters)
		if err != nil {
			return 0, getBuilderFuncCRUDError("select count", err)
		}
	}

	row := c.db.QueryRowContext(ctx, query, sqlfilters.FiltersInterfaces(options.Filters)...)
	var cnt int64
	err = row.Scan(&cnt)
	if err != nil {
		return 0, getDBFuncCRUDError("query row scan", err)
	}

	return cnt, nil
}
