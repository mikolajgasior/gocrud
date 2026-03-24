package structinfo

type FieldInfo struct {
	// Name contains name of the field.
	Name string
	// OverrideColumnName is the name of the column in the database table.
	OverrideColumnName string
	// OverrideColumnType is the type of the column in the database table.
	OverrideColumnType string
	// Default contains default value for the field.
	Default string
	// Unique is set to true if the field is marked with a tag 'uniq'.
	Unique bool
	// Password is set to true if the field is marked with a "password" tag.
	Password bool
	// Ignored is set to true if the field is marked with a "-" tag.
	Ignored bool
}
