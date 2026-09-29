package tablesync

import (
	"strings"
)

type TableData struct {
	Headers []string
	Rows    [][]string
}

func IsSimilarTable(baseHeaders, newHeaders []string) bool {
	if len(newHeaders) == 0 {
		return false
	}
	similarity := float64(MatchingHeaderCount(baseHeaders, newHeaders)) / float64(len(newHeaders))
	return similarity >= 0.4
}

func MatchingHeaderCount(baseHeaders, newHeaders []string) int {
	if len(baseHeaders) == 0 || len(newHeaders) == 0 {
		return 0
	}

	baseMap := make(map[string]bool)
	for _, header := range baseHeaders {
		cleanHeader := CleanString(header)
		if cleanHeader != "" {
			baseMap[cleanHeader] = true
		}
	}

	matches := 0
	for _, header := range newHeaders {
		if baseMap[CleanString(header)] {
			matches++
		}
	}
	return matches
}

func SelectMainTable(tables []TableData) (int, int) {
	if len(tables) == 0 {
		return -1, 0
	}

	mainIndex := 0
	maxMatches := -1
	for candidateIndex, candidate := range tables {
		matches := 0
		for otherIndex, other := range tables {
			if candidateIndex == otherIndex {
				continue
			}
			matches += MatchingHeaderCount(candidate.Headers, other.Headers)
		}
		if matches > maxMatches {
			mainIndex = candidateIndex
			maxMatches = matches
		}
	}

	return mainIndex, maxMatches
}

func SyncTwoTables(base TableData, addition TableData) TableData {
	headerIndexMap := make(map[string]int)
	for i, h := range base.Headers {
		headerIndexMap[CleanString(h)] = i
	}

	for _, h := range addition.Headers {
		cleanH := CleanString(h)
		if _, exists := headerIndexMap[cleanH]; !exists {
			base.Headers = append(base.Headers, h)
			headerIndexMap[cleanH] = len(base.Headers) - 1

			for rIdx := range base.Rows {
				base.Rows[rIdx] = append(base.Rows[rIdx], "")
			}
		}
	}

	for _, addRow := range addition.Rows {
		newRow := make([]string, len(base.Headers))
		for colIdx, val := range addRow {
			if colIdx >= len(addition.Headers) {
				continue
			}
			colName := CleanString(addition.Headers[colIdx])
			if targetColIdx, ok := headerIndexMap[colName]; ok {
				newRow[targetColIdx] = strings.TrimSpace(val)
			}
		}

		if existingRowIdx, exists := findMatchingRow(base.Rows, newRow); exists {
			mergeRows(base.Rows[existingRowIdx], newRow)
			continue
		}

		if !isRowEmpty(newRow) {
			base.Rows = append(base.Rows, newRow)
		}
	}

	return base
}

func findMatchingRow(rows [][]string, candidate []string) (int, bool) {
	for rowIndex, row := range rows {
		matchingCells := 0
		conflictingCells := false
		for colIndex := 0; colIndex < len(candidate); colIndex++ {
			incoming := CleanString(candidate[colIndex])
			if incoming == "" || colIndex >= len(row) {
				continue
			}
			current := CleanString(row[colIndex])
			if current == "" {
				continue
			}
			if current != incoming {
				conflictingCells = true
				break
			}
			matchingCells++
		}
		if matchingCells > 0 && !conflictingCells {
			return rowIndex, true
		}
	}
	return -1, false
}

func mergeRows(target, incoming []string) {
	for colIndex, value := range incoming {
		if colIndex >= len(target) || strings.TrimSpace(value) == "" {
			continue
		}
		if strings.TrimSpace(target[colIndex]) == "" {
			target[colIndex] = strings.TrimSpace(value)
		}
	}
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
