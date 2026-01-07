package crud

import (
	"context"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

type GetOptions struct {
	Order               []string
	Limit               int
	Offset              int
	Filters             *sqlbuilder.Filters
	RowObjTransformFunc func(interface{}) interface{}
}

func (c *CRUD) Get(ctx context.Context, newObjFunc func() interface{}, options GetOptions) ([]interface{}, error) {
	obj := newObjFunc()

	builder, err := c.builder(obj)
	if err != nil {
		return nil, getBuilderObjectCRUDError(err)
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
			return nil, getObjFuncCRUDError("select", err)
		}
	} else {
		query, err = builder.Select(options.Order, options.Limit, options.Offset, options.Filters)
		if err != nil {
			return nil, getBuilderFuncCRUDError("select", err)
		}
	}

	rows, err := c.db.QueryContext(ctx, query, sqlbuilder.FiltersInterfaces(options.Filters)...)
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
