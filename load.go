package crud

import (
	"context"
	"database/sql"
	"strconv"
)

type LoadOptions struct {
	Unused bool
}

func (c *CRUD) Load(ctx context.Context, obj interface{}, id string, options LoadOptions) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return getConvertIDToIntCRUDError(err)
	}

	builder, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}

	var query string

	// if the object has a SelectByID method, use it.
	if selectByIDerImpl, ok := obj.(selectByIDQueryBuilder); ok {
		query, err = selectByIDerImpl.SelectByIDQuery()
		if err != nil {
			return getObjFuncCRUDError("select by id query", err)
		}
	} else {
		query = builder.SelectByID()
	}

	err = c.db.QueryRowContext(ctx, query, int64(idInt)).Scan(ObjFieldInterfaces(obj, true)...)
	switch {
	case err == sql.ErrNoRows:
		ZeroObjFields(obj)

		return nil

	case err != nil:
		return getDBFuncCRUDError("query row", err)

	default:
		return nil
	}
}
