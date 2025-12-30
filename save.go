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
		return ErrCRUD{
			Op:  "o.builder",
			Err: err,
		}
	}

	ok, violations, err := Validate(obj, nil, c.tagName)
	if err != nil {
		return ErrCRUD{
			Op:  "Validate",
			Err: err,
		}
	}

	if !ok {
		return ErrCRUD{
			Op: "Validate",
			Err: &ErrValidation{
				Violations: violations,
			},
		}
	}

	objID := ObjIDValue(obj)

	// populate created and modified columns
	if options.ModifiedAt != 0 && options.ModifiedBy != 0 && builder.HasModificationFields() {
		SetObjModified(obj, options.ModifiedAt, options.ModifiedBy)
		if objID == 0 {
			SetObjCreated(obj, options.ModifiedAt, options.ModifiedBy)
		}
	}

	// encrypt all password fields
	passFields := builder.PasswordFields()
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
		uniqFields := builder.UniqueFields()
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
				return ErrCRUD{
					Op:  "o.GetCount",
					Err: err,
				}
			}

			if count > 0 {
				return ErrUniq{
					Field: uniqField,
				}
			}
		}
	}

	objIDInterface := ObjIDInterface(obj)

	// update
	if objID != 0 {
		// do no try to insert if NoInsert is set
		// TODO: error handling, we should check if object exists - for now nothing happens, UPDATE gets executed and updates nothing
		if options.NoInsert {
			_, err = c.db.ExecContext(ctx, builder.UpdateByID(), append(ObjFieldInterfaces(obj, false), objIDInterface)...)
		} else {
			// try to insert - if ID already exists then try to update it
			_, err = c.db.ExecContext(ctx, builder.InsertOnConflictUpdate(), append(ObjFieldInterfaces(obj, true), ObjFieldInterfaces(obj, false)...)...)
		}

		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) {
				// 23505 = unique_violation
				if pgErr.Code == "23505" {
					// You can also look at pgErr.Constraint to know which
					// unique index triggered the error (if you need that detail).
					return ErrUniq{
						Field: pgErr.Constraint,
					}
				}
			}

			return ErrCRUD{
				Op:  "o.db.QueryRow",
				Err: err,
			}
		}

		return nil
	}

	// insert
	err = c.db.QueryRowContext(ctx, builder.Insert(), ObjFieldInterfaces(obj, false)...).Scan(objIDInterface)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			// 23505 = unique_violation
			if pgErr.Code == "23505" {
				// You can also look at pgErr.Constraint to know which
				// unique index triggered the error (if you need that detail).
				return ErrUniq{
					Field: pgErr.Constraint,
				}
			}
		}

		return ErrCRUD{
			Op:  "o.db.QueryRow",
			Err: err,
		}
	}

	return nil

}
