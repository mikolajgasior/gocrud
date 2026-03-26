package builder

import (
	"codeberg.org/mikolajgasior/gocrud/pkg/filters"
	postgresQueryBuilder "codeberg.org/mikolajgasior/gocrud/pkg/sql/builder/querybuilder/postgres"
	postgresQueryContainer "codeberg.org/mikolajgasior/gocrud/pkg/sql/builder/querycontainer/postgres"
	"codeberg.org/mikolajgasior/gocrud/pkg/structinfo"
)

const (
	DefaultTagName = "sql"
)

// Builder reflects the object to generate and cache PostgreSQL queries (CREATE TABLE, INSERT, UPDATE, etc.).
// Database table and column names are lowercase with an underscore, and they are generated from field names.
type Builder struct {
	tagName string

	structInfo     *structinfo.StructInfo
	queryContainer *postgresQueryContainer.QueryContainer
	queryBuilder   *postgresQueryBuilder.QueryBuilder
}

// New takes a struct and returns a Builder instance.
func New(obj interface{}, options Options) *Builder {
	builder := &Builder{}

	builder.tagName = DefaultTagName
	if options.TagName != "" {
		builder.tagName = options.TagName
	}

	builder.structInfo = structinfo.New(obj, builder.tagName)
	builder.queryContainer = postgresQueryContainer.New(obj, builder.structInfo, options.TableNamePrefix)
	builder.queryBuilder = postgresQueryBuilder.New(builder.structInfo, builder.queryContainer)

	return builder
}

// DropTable returns an SQL query for dropping the table.
func (b *Builder) DropTable() string {
	return b.queryContainer.DropTable + ";"
}

// CreateTable returns an SQL query for creating the table.
func (b *Builder) CreateTable() string {
	return b.queryContainer.CreateTable + ";"
}

// Insert returns an SQL query for inserting a new object to the table.
func (b *Builder) Insert() string {
	return b.queryContainer.Insert + ";"
}

// UpdateByID returns an SQL query for updating an object by their ID.
func (b *Builder) UpdateByID() string {
	return b.queryContainer.UpdateByID + ";"
}

// InsertOnConflictUpdate returns an SQL query for inserting when a conflict is detected.
func (b *Builder) InsertOnConflictUpdate() string {
	return b.queryContainer.InsertOnConflictUpdate + ";"
}

// SelectByID returns an SQL query for selecting an object by its ID.
func (b *Builder) SelectByID() string {
	return b.queryContainer.SelectByID + ";"
}

// DeleteByID returns an SQL query for deleting an object by its ID.
func (b *Builder) DeleteByID() string {
	return b.queryContainer.DeleteByID + ";"
}

// DatabaseColumnToFieldName takes a database column and converts it to a struct field name.
func (b *Builder) DatabaseColumnToFieldName(n string) string {
	return b.queryContainer.ColumnFieldName[n]
}

// HasModificationFields returns true if all the following int64 fields are present: CreatedAt, CreatedBy, ModifiedAt, ModifiedBy.
func (b *Builder) HasModificationFields() bool {
	return b.structInfo.ModificationFields
}

// HasAliasedColumnNames returns true if any of the fields have column names in the format of "alias.column_name".
func (b *Builder) HasAliasedColumnNames() bool {
	return b.structInfo.AliasedColumnNames
}

// UniqueFields returns a list with field names that are unique.
func (b *Builder) UniqueFields() []string {
	return b.structInfo.UniqueFields
}

// PasswordFields returns a list with field names that are passwords.
func (b *Builder) PasswordFields() []string {
	return b.structInfo.PasswordFields
}

// Select returns a SELECT query with a WHERE condition built from 'filters' (field-value pairs).
// Struct fields in the 'filters' argument are sorted alphabetically. Hence, when used with a database connection, their values (or pointers to it) must be sorted as well.
// Columns in the SELECT query are ordered the same way as they are defined in the struct: SELECT field1_column, field2_column, ... etc.
func (b *Builder) Select(order []string, limit int, offset int, filters *filters.Filters) (string, error) {
	query, err := b.queryBuilder.Select(order, limit, offset, filters)
	if err != nil {
		return "", getQueryBuilderError("Select", err)
	}

	return query + ";", nil
}

// SelectCount returns a SELECT COUNT(*) query to count rows with a WHERE condition built from 'filters' (field-value pairs).
// Struct fields in the 'filters' argument are sorted alphabetically. Hence, when used with a database connection, their values (or pointers to it) must be sorted as well.
func (b *Builder) SelectCount(filters *filters.Filters) (string, error) {
	query, err := b.queryBuilder.SelectCount(filters)
	if err != nil {
		return "", getQueryBuilderError("SelectCount", err)
	}

	return query + ";", nil
}

// Delete returns a DELETE query with a WHERE condition built from 'filters' (field-value pairs).
// Struct fields in the 'filters' argument are sorted alphabetically. Hence, when used with a database connection, their values (or pointers to it) must be sorted as well.
func (b *Builder) Delete(filters *filters.Filters) (string, error) {
	query, err := b.queryBuilder.Delete(filters)
	if err != nil {
		return "", getQueryBuilderError("Delete", err)
	}

	return query + ";", nil
}

// DeleteReturningID returns a DELETE query with a WHERE condition built from 'filters' (field-value pairs) with RETURNING id.
// Struct fields in the 'filters' argument are sorted alphabetically. Hence, when used with a database connection, their values (or pointers to it) must be sorted as well.
func (b *Builder) DeleteReturningID(filters *filters.Filters) (string, error) {
	query, err := b.queryBuilder.DeleteReturningID(filters)
	if err != nil {
		return "", getQueryBuilderError("DeleteReturningID", err)
	}

	return query + ";", nil
}

// Update returns an UPDATE query where specified struct fields (columns) are updated and rows match specific WHERE condition built from 'filters' (field-value pairs).
// Struct fields in 'values' and the 'filters' arguments are sorted alphabetically. Hence, when used with a database connection, their values (or pointers to it) must be sorted as well.
func (b *Builder) Update(values map[string]interface{}, filters *filters.Filters) (string, error) {
	query, err := b.queryBuilder.Update(values, filters)
	if err != nil {
		return "", getQueryBuilderError("Update", err)
	}

	return query + ";", nil
}
