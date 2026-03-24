package crud

import (
	"context"

	sqlfilters "miko.gs/gocrud/pkg/filters"
)

type DeleteMultipleOptions struct {
	Filters            *sqlfilters.Filters
	CascadeDeleteDepth int8
}

// DeleteMultiple removes objects from the database based on specified filters
func (c *CRUD) DeleteMultiple(ctx context.Context, obj interface{}, options DeleteMultipleOptions) error {
	builder, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}

	err = ValidateFilters(obj, options.Filters, c.tagName)
	if err != nil {
		return err
	}

	var query string
	// if the object has a DeleteReturningID method, use it.
	if deleteReturningIDerImpl, ok := obj.(deleteReturningIDQueryBuilder); ok {
		query, err = deleteReturningIDerImpl.DeleteReturningIDQuery()
		if err != nil {
			return getObjFuncCRUDError("delete returning id query", err)
		}
	} else {
		query, err = builder.DeleteReturningID(options.Filters)
		if err != nil {
			return getBuilderFuncCRUDError("delete returning id", err)
		}
	}

	rows, err := c.db.QueryContext(ctx, query, sqlfilters.FiltersInterfaces(options.Filters)...)
	if err != nil {
		return getDBFuncCRUDError("query", err)
	}
	defer rows.Close()

	returnedIDs := []int64{}

	for rows.Next() {
		var returnedID int64
		errScan := rows.Scan(&returnedID)
		if errScan != nil {
			return getDBFuncCRUDError("rows scan", errScan)
		}

		returnedIDs = append(returnedIDs, returnedID)
	}

	if options.CascadeDeleteDepth < 3 {
		// Loop through the fields to cascade-delete.
		errCascadeDelete := c.runOnDelete(ctx, obj, returnedIDs, options.CascadeDeleteDepth)
		if errCascadeDelete != nil {
			return getCascadingDeleteCRUDError(errCascadeDelete)
		}
	}

	return nil
}
