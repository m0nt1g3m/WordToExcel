package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"wordtoexcel/internal/docxparser"
	"wordtoexcel/internal/excelwriter"

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

		sheetName := excelwriter.CleanSheetName(filepath.Base(docPath))
		if idx == 0 {
			if err := excel.SetSheetName("Sheet1", sheetName); err != nil {
				return "", fmt.Errorf("creating a sheet for file %s: %w", filepath.Base(docPath), err)
			}
		} else {
			// Each next document creates a separate sheet so existing data is not overwritten.
			sheetName = uniqueSheetName(excel, sheetName)
			if _, err := excel.NewSheet(sheetName); err != nil {
				return "", fmt.Errorf("creating a sheet for file %s: %w", filepath.Base(docPath), err)
			}
		}

		// Extract all tables found.
		tables, err := docxparser.ExtractTablesFromDocx(docPath)
		if err != nil {
			status(fmt.Sprintf("Error reading file %s: %v", filepath.Base(docPath), err))
			if progress != nil {
				progress(idx+1, len(*fileList))
			}
			continue
		}

		if len(tables) == 0 {
			status(fmt.Sprintf("No tables found in file %s", filepath.Base(docPath)))
			if progress != nil {
				progress(idx+1, len(*fileList))
			}
			continue
		}

		// nextRow stores the first free row on the sheet.
		// After each written table, the row position moves downward.
		nextRow := 1
		for tableIndex, table := range tables {
			nextRow = excelwriter.WriteTableToSheet(excel, sheetName, table, nextRow)
			status(fmt.Sprintf("File %s: table %d placed on sheet %s", filepath.Base(docPath), tableIndex+1, sheetName))
		}
		if progress != nil {
			progress(idx+1, len(*fileList))
		}
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("determining the application startup folder: %w", err)
	}
	tempFile, err := os.CreateTemp(workingDir, ".word-to-excel-*.xlsx")
	if err != nil {
		return "", fmt.Errorf("creating temporary file: %w", err)
	}
	tempPath := tempFile.Name()
	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("closing temporary file: %w", err)
	}
	if err := hideTempFile(tempPath); err != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("hiding temporary file: %w", err)
	}
	if err := excel.SaveAs(tempPath); err != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("saving Excel file: %w", err)
	}

	return tempPath, nil
}

// uniqueSheetName creates a new unique sheet for the document if the name is already taken.
func uniqueSheetName(excel *excelize.File, baseName string) string {
	existing := make(map[string]struct{})
	for _, name := range excel.GetSheetList() {
		existing[strings.ToLower(name)] = struct{}{}
	}
	if _, exists := existing[strings.ToLower(baseName)]; !exists {
		return baseName
	}

	for suffixIndex := 2; ; suffixIndex++ {
		suffix := fmt.Sprintf(" (%d)", suffixIndex)
		baseRunes := []rune(baseName)
		maxBaseLength := 31 - len([]rune(suffix))
		if len(baseRunes) > maxBaseLength {
			baseRunes = baseRunes[:maxBaseLength]
		}
		candidate := string(baseRunes) + suffix
		if _, exists := existing[strings.ToLower(candidate)]; !exists {
			return candidate
		}
	}
}
