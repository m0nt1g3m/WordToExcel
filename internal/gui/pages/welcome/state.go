package welcome

import (
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/manager"

	"gioui.org/widget"
)

type MainPage struct {
	SourceButton  *widget.Clickable
	FolderButton  *widget.Clickable
	ConvertButton *widget.Clickable
	manager       *manager.Manager
	status        *components.Status
}

func NewMainPage(mgr *manager.Manager, statusComponent *components.Status) *MainPage {
	return &MainPage{
		SourceButton:  &widget.Clickable{},
		FolderButton:  &widget.Clickable{},
		ConvertButton: &widget.Clickable{},
		manager:       mgr,
		status:        statusComponent,
	}
}
