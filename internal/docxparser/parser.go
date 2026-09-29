package docxparser

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/nguyenthenguyen/docx"

	"wordtoexcel/internal/tablesync"
)

func ExtractTablesFromDocx(path string) ([]tablesync.TableData, error) {
	doc, err := docx.ReadDocxFile(path)
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	content := doc.Editable().GetContent()
	var tables []tablesync.TableData

	tableBlocks := strings.Split(content, "<w:tbl>")
	for i, tbl := range tableBlocks {
		if i == 0 {
			continue
		}
		tblXML := strings.Split(tbl, "</w:tbl>")[0]
		parsedTable := parseXMLTable(tblXML)

		if len(parsedTable.Rows) > 0 {
			tables = append(tables, parsedTable)
		}
	}

	return tables, nil
}

func parseXMLTable(tblXML string) tablesync.TableData {
	var table tablesync.TableData
	rows, err := parseRowsFromXML(tblXML)
	if err != nil || len(rows) == 0 {
		return table
	}

	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	for i, row := range rows {
		if len(row) < maxCols {
			row = append(row, make([]string, maxCols-len(row))...)
		}
		if i == 0 {
			table.Headers = row
			continue
		}
		if !isRowEmpty(row) {
			table.Rows = append(table.Rows, row)
		}
	}

	return table
}

func parseRowsFromXML(tblXML string) ([][]string, error) {
	decoder := xml.NewDecoder(strings.NewReader(tblXML))
	var rows [][]string
	var currentRow []string
	var captureText strings.Builder
	var inCell bool
	var inText bool

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch v := tok.(type) {
		case xml.StartElement:
			switch v.Name.Local {
			case "tr":
				currentRow = nil
			case "tc":
				captureText.Reset()
				inCell = true
			case "t":
				inText = true
			}
		case xml.CharData:
			if inCell && inText {
				captureText.WriteString(string(v))
			}
		case xml.EndElement:
			switch v.Name.Local {
			case "t":
				inText = false
			case "tc":
				if inCell {
					currentRow = append(currentRow, strings.TrimSpace(captureText.String()))
					captureText.Reset()
					inCell = false
				}
			case "tr":
				if !isRowEmpty(currentRow) {
					rows = append(rows, currentRow)
				}
				currentRow = nil
			}
		}
	}

	return rows, nil
}

func isRowEmpty(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
