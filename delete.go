package crud

import "context"

type DeleteOptions struct{}

// Delete removes an object from the database table only when ID field is set (greater than 0).
// Once deleted from the DB, all field values are zeroed
func (c *CRUD) Delete(ctx context.Context, obj interface{}, options DeleteOptions) error {
	builder, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}

	id := ObjIDValue(obj)

	if id == 0 {
		return nil
	}

	var query string
	// if the object has a DeleteByID method, use it.
	if deleteByIDerImpl, ok := obj.(deleteByIDer); ok {
		query, err = deleteByIDerImpl.DeleteByID()
		if err != nil {
			return getObjFuncCRUDError("delete by id", err)
		}
	} else {
		query, err = builder.DeleteByID()
		if err != nil {
			return getBuilderFuncCRUDError("delete by id", err)
		}
	}

	_, err = c.db.ExecContext(ctx, query, ObjIDInterface(obj))
	if err != nil {
		return getDBFuncCRUDError("exec", err)
	}
	ZeroObjFields(obj)

	// Loop through fields and cascade-delete.
	err = c.runOnDelete(ctx, obj, []int64{id}, 0)
	if err != nil {
		return getCascadingDeleteCRUDError(err)
	}

	return nil
}
