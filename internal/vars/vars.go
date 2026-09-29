package vars

type AppState int

const (
	StateInit AppState = iota
	StateFilesSelected
	StateProcessing
	StateReady
	StateError
)
