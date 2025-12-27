package crud

import (
	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

type DeleteMultipleOptions struct {
	Filters            *sqlbuilder.Filters
	CascadeDeleteDepth int8
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

	query, err := builder.DeleteReturningID(options.Filters)
	if err != nil {
		return ErrCRUD{
			Op:  "builder.Delete",
			Err: err,
		}
	}

	rows, err := c.db.Query(query, sqlbuilder.FiltersInterfaces(options.Filters)...)
	if err != nil {
		return ErrCRUD{
			Op:  "o.db.Query",
			Err: err,
		}
	}
	defer rows.Close()

	returnedIDs := []int64{}

	for rows.Next() {
		var returnedID int64
		errScan := rows.Scan(&returnedID)
		if errScan != nil {
			return ErrCRUD{
				Op:  "db.Query",
				Err: errScan,
			}
		}

		returnedIDs = append(returnedIDs, returnedID)
	}

	if options.CascadeDeleteDepth < 3 {
		// Loop through the fields to cascade-delete.
		errCascadeDelete := c.runOnDelete(obj, returnedIDs, options.CascadeDeleteDepth)
		if errCascadeDelete != nil {
			return ErrCRUD{
				Op:  "o.runOnDelete",
				Err: errCascadeDelete,
			}
		}
	}

	return nil
}
