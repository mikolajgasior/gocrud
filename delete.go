package crud

type DeleteOptions struct{}

// Delete removes an object from the database table only when ID field is set (greater than 0).
// Once deleted from the DB, all field values are zeroed
func (c *CRUD) Delete(obj interface{}, options DeleteOptions) error {
	builder, err := c.builder(obj)
	if err != nil {
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	id := ObjIDValue(obj)

	if id == 0 {
		return nil
	}
	_, err = c.db.Exec(builder.DeleteByID(), ObjIDInterface(obj))
	if err != nil {
		return ErrCRUD{
			Op:  "o.db.Exec",
			Err: err,
		}
	}
	ZeroObjFields(obj)

	// Loop through fields and cascade-delete.
	err = c.runOnDelete(obj, []int64{id}, 0)
	if err != nil {
		return ErrCRUD{
			Op:  "o.runOnDelete",
			Err: err,
		}
	}

	return nil
}
