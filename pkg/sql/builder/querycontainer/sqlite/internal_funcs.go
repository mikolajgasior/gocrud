package sqlite

import "fmt"

func columnDefinitionFromField(fieldName string, fieldType string, isUnique bool, overrideColumnType string) string {
	if fieldName == "ID" {
		return "INTEGER PRIMARY KEY"
	}

	if overrideColumnType != "" {
		definition := fmt.Sprintf("%s NOT NULL DEFAULT ''", overrideColumnType)
		if isUnique {
			definition += " UNIQUE"
		}
		return definition
	}

	var definition string
	switch fieldType {
	case "string":
		definition = "TEXT NOT NULL DEFAULT ''"
	case "bool":
		definition = "INTEGER NOT NULL DEFAULT 0"
	case "float32", "float64":
		definition = "REAL NOT NULL DEFAULT 0"
	default:
		// all int/uint variants
		definition = "INTEGER NOT NULL DEFAULT 0"
	}

	if isUnique {
		definition += " UNIQUE"
	}

	return definition
}
