package manager

import (
	"fmt"
	"log"
	"os"
	"sync"
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/utils"
	"wordtoexcel/internal/vars"

	"gioui.org/app"
)

type Manager struct {
	Window         *app.Window
	AppState       vars.AppState
	fileList       []string
	outputFilePath string

	statusMessages []string // list of messages displayed in the interface console
	processing     bool     // flag indicating an active conversion process
	readyToSave    bool     // flag allowing the save button to be displayed
	progressDone   int      // number of successfully processed files
	progressTotal  int      // total number of files in the current task
	tempFilePath   string   // path to the temporary result before the user saves it

	mu sync.RWMutex
}

// State returns the current application state.
func (m *Manager) State() vars.AppState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.AppState
}

func (m *Manager) SetState(state vars.AppState) {
	m.mu.Lock()
	m.AppState = state
	m.mu.Unlock()

	if m.Window != nil {
		m.Window.Invalidate()
	}
}

func (m *Manager) SelectedFilesCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.fileList)
}

func (m *Manager) GetStatus() ([]string, string, bool, bool, int, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]string(nil), m.statusMessages...), m.tempFilePath, m.processing, m.readyToSave, m.progressDone, m.progressTotal
}

func (m *Manager) AddStatusMessage(msg string) {
	m.mu.Lock()
	m.statusMessages = append(m.statusMessages, msg)
	m.mu.Unlock()
	if m.Window != nil {
		m.Window.Invalidate()
	}
}

func (m *Manager) SetProgress(done, total int) {
	m.mu.Lock()
	m.progressDone = done
	m.progressTotal = total
	m.mu.Unlock()
	if m.Window != nil {
		m.Window.Invalidate()
	}
}

// SelectFiles opens the Word file selection dialog.
func (m *Manager) SelectFiles() {
	components.SelectFiles(&m.fileList, m.AddStatusMessage)
	if m.SelectedFilesCount() > 0 {
		m.SetState(vars.StateFilesSelected)
	}
}

// SelectFolder opens a folder picker and searches it for Word documents.
func (m *Manager) SelectFolder() {
	components.SelectFolder(&m.fileList, m.AddStatusMessage)
	if m.SelectedFilesCount() > 0 {
		m.SetState(vars.StateFilesSelected)
	}
}

// ConvertFiles starts the document conversion process.
func (m *Manager) ConvertFiles() {
	m.mu.Lock()
	oldTempFilePath := m.tempFilePath
	m.statusMessages = nil
	m.processing = true
	m.readyToSave = false
	m.tempFilePath = ""
	m.progressDone = 0
	m.progressTotal = 0
	m.mu.Unlock()
	removeTempFile(oldTempFilePath)

	m.SetState(vars.StateProcessing)

	go func() {
		tempPath, err := utils.ConvertWordToExcel(&m.fileList, m.AddStatusMessage, m.SetProgress)

		m.mu.Lock()
		m.processing = false
		if err == nil {
			m.tempFilePath = tempPath
			m.readyToSave = true
		}
		m.mu.Unlock()

		if err != nil {
			m.SetState(vars.StateError)
			m.AddStatusMessage(fmt.Sprintf("Error: %v", err))
		} else {
			m.SetState(vars.StateReady)
			m.AddStatusMessage("All files processed.")
		}
	}()
}

func (m *Manager) SaveResult() {
	go func() {
		m.mu.RLock()
		tempPath := m.tempFilePath
		m.mu.RUnlock()

		if err := components.SaveFile(tempPath, &m.outputFilePath); err != nil {
			m.AddStatusMessage(fmt.Sprintf("Save error: %v", err))
			return
		}

		m.mu.Lock()
		m.readyToSave = false
		m.tempFilePath = ""
		m.mu.Unlock()
		removeTempFile(tempPath)

		m.AddStatusMessage(fmt.Sprintf("File saved: %s", m.outputFilePath))
	}()
}

func (m *Manager) Reset() {
	m.mu.Lock()
	tempFilePath := m.tempFilePath
	m.AppState = vars.StateInit
	m.fileList = nil
	m.outputFilePath = ""
	m.statusMessages = nil
	m.processing = false
	m.readyToSave = false
	m.progressDone = 0
	m.progressTotal = 0
	m.tempFilePath = ""
	m.mu.Unlock()

	if tempFilePath != "" {
		removeTempFile(tempFilePath)
	}
	if m.Window != nil {
		m.Window.Invalidate()
	}
}

// Cleanup removes the temporary conversion result when the application exits.
func (m *Manager) Cleanup() {
	m.mu.Lock()
	tempFilePath := m.tempFilePath
	m.tempFilePath = ""
	m.readyToSave = false
	m.mu.Unlock()
	removeTempFile(tempFilePath)
}

func removeTempFile(path string) {
	if path == "" {
		return
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove temporary file %q: %v", path, err)
	}
}

func NewManager() *Manager {
	return &Manager{
		Window:         new(app.Window),
		AppState:       vars.StateInit,
		fileList:       make([]string, 0),
		outputFilePath: "",
	}
}
