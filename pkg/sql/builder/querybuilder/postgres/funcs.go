package postgres

import "fmt"

// LimitOffset returns a LIMIT and OFFSET clause for a query.
func LimitOffset(limit int, offset int) string {
	if limit == 0 {
		return ""
	}

	if offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}

	return fmt.Sprintf("LIMIT %d", limit)
}
