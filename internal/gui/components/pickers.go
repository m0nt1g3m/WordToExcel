package components

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
)

func SelectFiles(fileList *[]string, status func(string)) {
	if fileList == nil {
		return
	}

	files, err := zenity.SelectFileMultiple(
		zenity.Title("Select Word files"),
		zenity.FileFilters{
			{
				Name:     "Word Documents (*.docx, *.doc)",
				Patterns: []string{"*.docx", "*.doc"},
			},
		},
	)

	if err != nil {
		log.Println("Error selecting files:", err)
		return
	}

	if len(files) == 0 {
		log.Println("No files selected.")
		return
	}
	if len(*fileList) != 0 {
		clear(*fileList)
		*fileList = (*fileList)[:0]
	}
	for _, file := range files {
		if !strings.HasPrefix(filepath.Base(file), "~$") {
			println(fmt.Sprintf("Selected file: %s", file))
			if status != nil {
				status(fmt.Sprintf("Word file selected: %s", file))
			}
			*fileList = append(*fileList, file)
		}
	}
}

func SelectSource(fileList *[]string, status func(string), onSelected func(int)) {
	if fileList == nil {
		return
	}

	err := zenity.Question(
		"Select Word documents source",
		zenity.Title("Select source"),
		zenity.OKLabel("Select files"),
		zenity.ExtraButton("Select folder"),
		zenity.CancelLabel("Cancel"),
	)

	switch err {
	case nil:
		SelectFiles(fileList, status)
	case zenity.ErrExtraButton:
		SelectFolder(fileList, status)
	case zenity.ErrCanceled:
		return
	default:
		log.Println("Error selecting source type:", err)
	}

	log.Printf("Lenght of fileList: %d", len(*fileList))
	if onSelected != nil {
		onSelected(len(*fileList))
	}
}

func SelectFolder(fileList *[]string, status func(string)) {
	if fileList == nil {
		return
	}

	dirPath, err := zenity.SelectFile(
		zenity.Title("Select folder"),
		zenity.Directory(),
	)

	if err != nil {
		log.Println("Error selecting folder:", err)
		return
	}

	if dirPath == "" {
		log.Println("No folder selected.")
		return
	}

	var foundFiles []string

	err = filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".docx" || ext == ".doc" {
				if !strings.HasPrefix(filepath.Base(path), "~$") {
					println(fmt.Sprintf("Founded file: %s", path))
					if status != nil {
						status(fmt.Sprintf("Word file found: %s", path))
					}
					foundFiles = append(foundFiles, path)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Println("Error walking through directory:", err)
		return
	}

	if len(foundFiles) == 0 {
		log.Println("No Word files found in the selected folder.")
		return
	}
	if len(*fileList) != 0 {
		clear(*fileList)
		*fileList = (*fileList)[:0]
	}

	*fileList = append(*fileList, foundFiles...)

}

func SaveFile(sourcePath string, filePath *string) error {
	if sourcePath == "" || filePath == nil {
		return fmt.Errorf("file is not ready for saving")
	}

	path, err := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename("word_to_excel.xlsx"),
		zenity.FileFilters{
			{
				Name:     "Excel Files (*.xlsx, *.xls)",
				Patterns: []string{"*.xlsx", "*.xls"},
			},
		},
	)

	if err != nil {
		return err
	}

	*filePath = path

	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err = io.Copy(destination, source); err != nil {
		destination.Close()
		return err
	}
	if err = destination.Close(); err != nil {
		return err
	}

	return os.Remove(sourcePath)
}
