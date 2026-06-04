package gocrud

import (
	"context"

	sqlfilters "codeberg.org/mikolajgasior/gocrud/pkg/filters"
	sqlbuilder "codeberg.org/mikolajgasior/gocrud/pkg/sql/builder"
	validator "github.com/mikolajgasior/struct-validator"
)

type GetOptions struct {
	Order                    []string
	Limit                    int
	Offset                   int
	Filters                  *sqlfilters.Filters
	RowObjTransformFunc      func(interface{}) interface{}
	ConvertFiltersFromString bool
}

func (c *CRUD) Get(ctx context.Context, newObjFunc func() interface{}, options GetOptions) ([]interface{}, error) {
	obj := newObjFunc()

	builder, err := c.builder(obj)
	if err != nil {
		return nil, getBuilderObjectCRUDError(err)
	}

	// Filter values can be passed as string. We do not want any use of reflect outside of CRUD.
	if options.ConvertFiltersFromString && options.Filters != nil {
		newFilters := &sqlfilters.Filters{}

		for filterName, filterOpVal := range *(options.Filters) {
			// ignore the raw filters entirely, that's too complicated
			if filterName == sqlfilters.Raw {
				continue
			}

			valueAsString, ok := filterOpVal.Val.(string)
			if !ok {
				return nil, getObjInvalidCRUDError(map[string]uint64{
					filterName: validator.FailType,
				})
			}

			ok, valueAsFieldType := sqlbuilder.StructFieldValueFromString(obj, filterName, valueAsString)
			if !ok {
				return nil, getObjInvalidCRUDError(map[string]uint64{
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
		return nil, err
	}

	var returnValue []interface{}

	var query string

	// if the object has a Select method, use it.
	if selectedImpl, ok := obj.(selectQueryBuilder); ok {
		query, err = selectedImpl.SelectQuery(options.Order, options.Limit, options.Offset, options.Filters)
		if err != nil {
			return nil, getObjFuncCRUDError("select query", err)
		}
	} else {
		query, err = builder.Select(options.Order, options.Limit, options.Offset, options.Filters)
		if err != nil {
			return nil, getBuilderFuncCRUDError("select", err)
		}
	}

	rows, err := c.db.QueryContext(ctx, query, sqlfilters.FiltersInterfaces(options.Filters)...)
	if err != nil {
		return nil, getDBFuncCRUDError("query", err)
	}
	defer rows.Close()

	for rows.Next() {
		newObj := newObjFunc()
		err = rows.Scan(ObjFieldInterfaces(newObj, true)...)
		if err != nil {
			return nil, getDBFuncCRUDError("rows scan", err)
		}

		// If options.RowObjTransformFunc is defined, then call it on the row.
		if options.RowObjTransformFunc != nil {
			returnValue = append(returnValue, options.RowObjTransformFunc(newObj))
			continue
		}

		// Normal append.
		returnValue = append(returnValue, newObj)
	}

	return returnValue, nil
}
