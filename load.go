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

	bldr, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}

	var query string

	// if the object has a SelectByID method, use it.
	if selectByIDerImpl, ok := obj.(selectByIDer); ok {
		query, err = selectByIDerImpl.SelectByID()
		if err != nil {
			return getObjFuncCRUDError("delete by id", err)
		}
	} else {
		query, err = bldr.SelectByID()
		if err != nil {
			return getBuilderFuncCRUDError("delete by id", err)
		}
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
