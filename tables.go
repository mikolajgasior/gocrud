package crud

// CreateTable creates a database table for the specified object.
func (c *CRUD) CreateTable(obj interface{}) error {
	builder, err := c.builder(obj)
	if err != nil {
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	_, err = c.db.Exec(builder.CreateTable())
	if err != nil {
		return ErrCRUD{
			Op:  "o.db.Exec",
			Err: err,
		}
	}

	return nil
}

// DropTable drops the database table for the specified object.
func (c *CRUD) DropTable(obj interface{}) error {
	builder, err := c.builder(obj)
	if err != nil {
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	_, err = c.db.Exec(builder.DropTable())
	if err != nil {
		return ErrCRUD{
			Op:  "o.db.Exec",
			Err: err,
		}
	}

	return nil
}
