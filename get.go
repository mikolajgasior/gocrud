package crud

import (
	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

type GetOptions struct {
	Order               []string
	Limit               int
	Offset              int
	Filters             *sqlbuilder.Filters
	RowObjTransformFunc func(interface{}) interface{}
}

func (c *CRUD) Get(newObjFunc func() interface{}, options GetOptions) ([]interface{}, error) {
	obj := newObjFunc()

	builder, err := c.builder(obj)
	if err != nil {
		return nil, ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	err = ValidateFilters(obj, options.Filters, c.tagName)
	if err != nil {
		return nil, err
	}

	var returnValue []interface{}
	query, err := builder.Select(options.Order, options.Limit, options.Offset, options.Filters)
	if err != nil {
		return nil, ErrCRUD{
			Op:  "builder.Select",
			Err: err,
		}
	}

	rows, err := c.db.Query(query, sqlbuilder.Interfaces(options.Filters)...)
	if err != nil {
		return nil, ErrCRUD{
			Op:  "o.db.Query",
			Err: err,
		}
	}
	defer rows.Close()

	for rows.Next() {
		newObj := newObjFunc()
		err = rows.Scan(ObjFieldInterfaces(newObj, true)...)
		if err != nil {
			return nil, ErrCRUD{
				Op:  "rows.Scan",
				Err: err,
			}
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
