package tablesync

import (
	"strings"
)

// TableData contains the table header row and all table rows.
type TableData struct {
	Headers []string
	Rows    [][]string
}

func isRowEmpty(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func CleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
