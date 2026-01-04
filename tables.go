package crud

import "context"

// CreateTable creates a database table for the specified object.
func (c *CRUD) CreateTable(ctx context.Context, obj interface{}) error {
	builder, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}

	var query string

	// if the object has a CreateTable method, use it.
	if createTableerImpl, ok := obj.(createTableer); ok {
		query, err = createTableerImpl.CreateTable()
		if err != nil {
			return getObjFuncCRUDError("create table", err)
		}
	} else {
		query, err = builder.CreateTable()
		if err != nil {
			return getBuilderFuncCRUDError("create table", err)
		}
	}

	_, err = c.db.ExecContext(ctx, query)
	if err != nil {
		return getDBFuncCRUDError("exec", err)
	}

	return nil
}

// DropTable drops the database table for the specified object.
func (c *CRUD) DropTable(ctx context.Context, obj interface{}) error {
	builder, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}

	var query string
	// if the object has a DropTable method, use it.
	if dropTableerImpl, ok := obj.(dropTableer); ok {
		query, err = dropTableerImpl.DropTable()
		if err != nil {
			return getObjFuncCRUDError("drop table", err)
		}
	} else {
		query, err = builder.DropTable()
		if err != nil {
			return getBuilderFuncCRUDError("drop table", err)
		}
	}

	_, err = c.db.ExecContext(ctx, query)
	if err != nil {
		return getDBFuncCRUDError("exec", err)
	}

	return nil
}
