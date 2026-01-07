package crud

import (
	"context"
	"errors"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type SaveOptions struct {
	NoInsert   bool
	ModifiedBy int64
	ModifiedAt int64
}

// Save takes an object, validates its field values, and saves it in the database.
// If ID is not present, then an INSERT will be performed.
// If ID is set, then an "upsert" is performed.
func (c *CRUD) Save(ctx context.Context, obj interface{}, options SaveOptions) error {
	builder, err := c.builder(obj)
	if err != nil {
		return getBuilderObjectCRUDError(err)
	}

	ok, violations, err := Validate(obj, nil, c.tagName)
	if err != nil {
		return getValidateObjCRUDError(err)
	}

	if !ok {
		return getObjInvalidCRUDError(violations)
	}

	objID := ObjIDValue(obj)

	// populate created and modified columns
	hasModificationFields, err := builder.HasModificationFields()
	if err != nil {
		return getBuilderFuncCRUDError("has modification fields", err)
	}

	if options.ModifiedAt != 0 && options.ModifiedBy != 0 && hasModificationFields {
		SetObjModified(obj, options.ModifiedAt, options.ModifiedBy)
		if objID == 0 {
			SetObjCreated(obj, options.ModifiedAt, options.ModifiedBy)
		}
	}

	// encrypt all password fields
	passFields, err := builder.PasswordFields()
	if err != nil {
		return getBuilderFuncCRUDError("password fields", err)
	}

	for _, passField := range passFields {
		fieldValue := ObjFieldValue(obj, passField)
		hash, err := bcrypt.GenerateFromPassword([]byte(fieldValue.(string)), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		SetObjStringField(obj, passField, string(hash))
	}

	// run GetCount on unique fields
	if c.flags&GetCountOnUniq > 0 {
		uniqFields, err := builder.UniqueFields()
		if err != nil {
			return getBuilderFuncCRUDError("unique fields", err)
		}
		for _, uniqField := range uniqFields {
			count, err := c.GetCount(ctx, obj,
				GetCountOptions{
					Filters: &sqlbuilder.Filters{
						uniqField: sqlbuilder.OpVal{
							Op:  sqlbuilder.OpEqual,
							Val: ObjFieldValue(obj, uniqField),
						},
						"ID": sqlbuilder.OpVal{
							Op:  sqlbuilder.OpNotEqual,
							Val: objID,
						},
					},
				},
			)
			if err != nil {
				return getCRUDFuncCRUDError("get count", err)
			}

			if count > 0 {
				return getUniqError(uniqField)
			}
		}
	}

	objIDInterface := ObjIDInterface(obj)

	// update
	if objID != 0 {
		// do no try to insert if NoInsert is set
		// TODO: error handling, we should check if object exists - for now nothing happens, UPDATE gets executed and updates nothing
		if options.NoInsert {
			var query string

			// if the object has an UpdateByID method, use it.
			if updateByIDerImpl, ok := obj.(updateByIDQueryBuilder); ok {
				query, err = updateByIDerImpl.UpdateByIDQuery()
				if err != nil {
					return getObjFuncCRUDError("delete by id", err)
				}
			} else {
				query, err = builder.UpdateByID()
				if err != nil {
					return getBuilderFuncCRUDError("update by id", err)
				}
			}
			_, err = c.db.ExecContext(ctx, query, append(ObjFieldInterfaces(obj, false), objIDInterface)...)
		} else {
			// try to insert - if ID already exists then try to update it

			var query string

			// if the object has an UpdateByID method, use it.
			if insertOnConflictUpdateerImpl, ok := obj.(insertOnConflictUpdateQueryBuilder); ok {
				query, err = insertOnConflictUpdateerImpl.InsertOnConflictUpdateQuery()
				if err != nil {
					return getObjFuncCRUDError("insert on conflict update", err)
				}
			} else {
				query, err = builder.InsertOnConflictUpdate()
				if err != nil {
					return getBuilderFuncCRUDError("insert on conflict update", err)
				}
			}
			_, err = c.db.ExecContext(ctx, query, append(ObjFieldInterfaces(obj, true), ObjFieldInterfaces(obj, false)...)...)
		}

		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) {
				// 23505 = unique_violation
				if pgErr.Code == "23505" {
					// You can also look at pgErr.Constraint to know which
					// unique index triggered the error (if you need that detail).
					return getUniqError(pgErr.Constraint)
				}
			}

			return getDBFuncCRUDError("query row", err)
		}

		return nil
	}

	// insert
	var query string
	// if the object has an Insert method, use it.
	if inserterImpl, ok := obj.(insertQueryBuilder); ok {
		query, err = inserterImpl.InsertQuery()
		if err != nil {
			return getObjFuncCRUDError("insert", err)
		}
	} else {
		query, err = builder.Insert()
		if err != nil {
			return getBuilderFuncCRUDError("insert", err)
		}
	}

	err = c.db.QueryRowContext(ctx, query, ObjFieldInterfaces(obj, false)...).Scan(objIDInterface)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			// 23505 = unique_violation
			if pgErr.Code == "23505" {
				// You can also look at pgErr.Constraint to know which
				// unique index triggered the error (if you need that detail).
				return getUniqError(pgErr.Constraint)
			}
		}

		return getDBFuncCRUDError("query row", err)
	}

	return nil

}
