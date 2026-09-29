package manager

import (
	"fmt"
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

	statusMessages []string // list of messages to display in the interface console
	processing     bool     // flag indicating an active conversion process
	readyToSave    bool     // flag allowing the display of the save button
	progressDone   int      // number of successfully processed files
	progressTotal  int      // total number of files in the current task
	tempFilePath   string   // path to the temporary result before saving by the user

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

// SelectFiles opens a dialog window for selecting files or folders.
// Upon successful selection, it changes the application state to StateFilesSelected.
func (m *Manager) SelectFiles() {
	components.SelectSource(&m.fileList, m.AddStatusMessage, func(count int) {
		if count > 0 {
			m.SetState(vars.StateFilesSelected)
		}
	})
}

// ConvertFiles starts the document conversion process.
// The method clears the previous log history and changes the application state to StateProcessing.
func (m *Manager) ConvertFiles() {
	m.mu.Lock()
	m.statusMessages = nil
	m.processing = true
	m.readyToSave = false
	m.tempFilePath = ""
	m.progressDone = 0
	m.progressTotal = 0
	m.mu.Unlock()

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
		m.mu.Unlock()

		m.AddStatusMessage(fmt.Sprintf("File saved: %s", m.outputFilePath))
	}()
}

func (m *Manager) Reset() {
	m.SetState(vars.StateInit)
}

func NewManager() *Manager {
	return &Manager{
		Window:         new(app.Window),
		AppState:       vars.StateInit,
		fileList:       make([]string, 0),
		outputFilePath: "",
	}
}
