package crud

import (
	"context"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

type GetCountOptions struct {
	Filters *sqlbuilder.Filters
}

// GetCount runs a 'SELECT COUNT(*)' query on the database with specified filters, order, limit and offset and returns count of rows
func (c *CRUD) GetCount(ctx context.Context, obj interface{}, options GetCountOptions) (int64, error) {
	builder, err := c.builder(obj)
	if err != nil {
		return 0, getBuilderObjectCRUDError(err)
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

	row := c.db.QueryRowContext(ctx, query, sqlbuilder.FiltersInterfaces(options.Filters)...)
	var cnt int64
	err = row.Scan(&cnt)
	if err != nil {
		return 0, getDBFuncCRUDError("query row scan", err)
	}

	return cnt, nil
}
