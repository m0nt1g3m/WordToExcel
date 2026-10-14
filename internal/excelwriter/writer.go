package excelwriter

import (
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"

	"wordtoexcel/internal/tablesync"
)

func WriteTableToSheet(f *excelize.File, sheet string, t tablesync.TableData, startRow int) int {
	for colIdx, h := range t.Headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, startRow)
		f.SetCellValue(sheet, cell, h)
	}

	for rowIdx, row := range t.Rows {
		for colIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, startRow+rowIdx+1)
			f.SetCellValue(sheet, cell, val)
		}
	}

	return startRow + len(t.Rows) + 2
}

func CleanSheetName(fileName string) string {
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	invalidChars := []string{"\\", "/", "?", "*", "[", "]"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}
	if name == "" {
		name = "Sheet"
	}
	runes := []rune(name)
	// The maximum number of characters in a sheet name is 31.
	if len(runes) > 31 {
		name = string(runes[:31])
	}
	return name
}
