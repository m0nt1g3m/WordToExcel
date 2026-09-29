package excelwriter

import (
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"

	"wordtoexcel/internal/tablesync"
)

func WriteTableToSheet(f *excelize.File, sheet string, t tablesync.TableData) {
	for colIdx, h := range t.Headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for rowIdx, row := range t.Rows {
		for colIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
		}
	}
}

func CleanSheetName(fileName string) string {
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	invalidChars := []string{"\\", "/", "?", "*", "[", "]"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}
	if len(name) > 31 {
		name = name[:31]
	}
	return name
}
