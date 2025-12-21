package crud

import (
	"database/sql"
	"strconv"
)

type LoadOptions struct {
	Unused bool
}

func (c *CRUD) Load(obj interface{}, id string, options LoadOptions) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return ErrCRUD{
			Op:  "atoi",
			Err: err,
		}
	}

	bldr, err := c.builder(obj)
	if err != nil {
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	err = c.db.QueryRow(bldr.SelectById(), int64(idInt)).Scan(ObjFieldInterfaces(obj, true)...)
	switch {
	case err == sql.ErrNoRows:
		ZeroObjFields(obj)

		return nil

	case err != nil:
		return ErrCRUD{
			Op:  "o.db.QueryRow",
			Err: err,
		}

	default:
		return nil
	}
}
