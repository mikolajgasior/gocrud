package gocrud

import "context"

// CreateTable creates a database table for the specified object.
func (c *CRUD) CreateTable(ctx context.Context, obj interface{}) error {
	builder, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}

	var query string

	// if the object has a CreateTable method, use it.
	if createTableerImpl, ok := obj.(createTableQueryBuilder); ok {
		query, err = createTableerImpl.CreateTableQuery()
		if err != nil {
			return getObjFuncCRUDError("create table", err)
		}
	} else {
		query = builder.CreateTable()
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
	if dropTableerImpl, ok := obj.(dropTableQueryBuilder); ok {
		query, err = dropTableerImpl.DropTableQuery()
		if err != nil {
			return getObjFuncCRUDError("drop table", err)
		}
	} else {
		query = builder.DropTable()
	}

	_, err = c.db.ExecContext(ctx, query)
	if err != nil {
		return getDBFuncCRUDError("exec", err)
	}

	return nil
}
