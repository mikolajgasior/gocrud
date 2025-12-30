package crud

import "context"

// CreateTable creates a database table for the specified object.
func (c *CRUD) CreateTable(ctx context.Context, obj interface{}) error {
	builder, err := c.builder(obj)
	if err != nil {
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	_, err = c.db.ExecContext(ctx, builder.CreateTable())
	if err != nil {
		return ErrCRUD{
			Op:  "o.db.Exec",
			Err: err,
		}
	}

	return nil
}

// DropTable drops the database table for the specified object.
func (c *CRUD) DropTable(ctx context.Context, obj interface{}) error {
	builder, err := c.builder(obj)
	if err != nil {
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	_, err = c.db.ExecContext(ctx, builder.DropTable())
	if err != nil {
		return ErrCRUD{
			Op:  "o.db.Exec",
			Err: err,
		}
	}

	return nil
}
