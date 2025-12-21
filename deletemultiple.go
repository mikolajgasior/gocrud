package crud

import (
	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

type DeleteMultipleOptions struct {
	Filters *sqlbuilder.Filters
}

// DeleteMultiple removes objects from the database based on specified filters
func (c *CRUD) DeleteMultiple(obj interface{}, options DeleteMultipleOptions) error {
	builder, err := c.builder(obj)
	if err != nil {
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	err = ValidateFilters(obj, options.Filters, c.tagName)
	if err != nil {
		return err
	}

	query, err := builder.Delete(options.Filters)
	if err != nil {
		return ErrCRUD{
			Op:  "builder.SelectCount",
			Err: err,
		}
	}

	_, err = c.db.Exec(query, sqlbuilder.Interfaces(options.Filters)...)
	if err != nil {
		return ErrCRUD{
			Op:  "o.db.Query",
			Err: err,
		}
	}

	return nil
}
