package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"wordtoexcel/internal/docxparser"
	"wordtoexcel/internal/excelwriter"
	"wordtoexcel/internal/tablesync"

	"github.com/xuri/excelize/v2"
)

func ConvertWordToExcel(fileList *[]string, status func(string), progress func(done, total int)) (string, error) {

	if fileList == nil || len(*fileList) == 0 {
		return "", fmt.Errorf("file list is empty")
	}

	excel := excelize.NewFile()
	defer excel.Close()

	for idx, docPath := range *fileList {
		status(fmt.Sprintf("Processing file %d of %d: %s", idx+1, len(*fileList), filepath.Base(docPath)))

		tables, err := docxparser.ExtractTablesFromDocx(docPath)
		if err != nil {
			status(fmt.Sprintf("Error reading file %s: %v", filepath.Base(docPath), err))
			continue
		}

		if len(tables) == 0 {
			status(fmt.Sprintf("No tables found in file %s", filepath.Base(docPath)))
			continue
		}

		sheetName := excelwriter.CleanSheetName(filepath.Base(docPath))

		if idx == 0 {
			excel.SetSheetName("Sheet1", sheetName)
		} else {
			// A new sheet is created for subsequent files
			excel.NewSheet(sheetName)
		}

		mainIndex, matches := tablesync.SelectMainTable(tables)
		status(fmt.Sprintf("File %s: main table #%d selected (column matches: %d)",
			filepath.Base(docPath), mainIndex+1, matches))

		syncedTable := tables[mainIndex]

		for i := range tables {
			if i == mainIndex {
				continue
			}
			currentTable := tables[i]

			if tablesync.IsSimilarTable(syncedTable.Headers, currentTable.Headers) {
				syncedTable = tablesync.SyncTwoTables(syncedTable, currentTable)
				status(fmt.Sprintf("File %s: table %d synced with the main one", filepath.Base(docPath), i+1))
			} else {
				extraSheetName := getExtraSheetName(sheetName, i+1)
				excel.NewSheet(extraSheetName)
				excelwriter.WriteTableToSheet(excel, extraSheetName, currentTable)
				status(fmt.Sprintf("File %s: table %d moved to sheet %s", filepath.Base(docPath), i+1, extraSheetName))
			}
		}

		// Table data is written strictly to the current sheetName
		excelwriter.WriteTableToSheet(excel, sheetName, syncedTable)
		if progress != nil {
			progress(idx+1, len(*fileList))
		}
	}

	tempFile, err := os.CreateTemp("", "word-to-excel-*.xlsx")
	if err != nil {
		return "", fmt.Errorf("creating temporary file: %w", err)
	}
	tempPath := tempFile.Name()
	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("closing temporary file: %w", err)
	}
	if err := excel.SaveAs(tempPath); err != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("saving Excel file: %w", err)
	}

	return tempPath, nil
}

func getExtraSheetName(baseName string, tableIndex int) string {
	suffix := fmt.Sprintf("_T%d", tableIndex)
	if len(baseName)+len(suffix) <= 31 {
		return baseName + suffix
	}
	prefixLen := 31 - len(suffix)
	if prefixLen > 0 {
		return baseName[:prefixLen] + suffix
	}
	return suffix
}
