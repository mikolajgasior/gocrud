package postgres

import (
	"fmt"
	"strings"
	"unicode"
)

// FieldToColumn converts struct field name to database column name.
func FieldToColumn(s string) string {
	if s == "ID" {
		return "id"
	}

	o := ""

	var prev rune
	for i, ch := range s {
		if i == 0 {
			o += strings.ToLower(string(ch))
			prev = ch
			continue
		}

		if unicode.IsUpper(ch) {
			if prev == 'I' && ch == 'D' {
				o += strings.ToLower(string(ch))
				continue
			}

			o += "_" + strings.ToLower(string(ch))
			prev = ch
			continue
		}

		o += string(ch)
		prev = ch
	}

	return o
}

// QuoteColumn quotes a column name.
func QuoteColumn(column string) string {
	// if column is in format of alias.name then quote only the name.
	if strings.Contains(column, ".") {
		parts := strings.SplitN(column, ".", 2)
		if len(parts) == 2 {
			return fmt.Sprintf(`%s."%s"`, parts[0], parts[1])
		}
	}

	return fmt.Sprintf(`"%s"`, column)
}
