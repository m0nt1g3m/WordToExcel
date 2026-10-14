package processing

import (
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/manager"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Page displays the interface while conversion is running.
type Page struct {
	manager *manager.Manager
	status  *components.Status
}

func NewPage(mgr *manager.Manager, statusComponent *components.Status) *Page {
	return &Page{
		manager: mgr,
		status:  statusComponent,
	}
}

func (p *Page) Draw(gtx layout.Context, theme *material.Theme) {
	// Get the current processing state and progress.
	statusMsgs, _, processing, _, done, total := p.manager.GetStatus()

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Render the top section of the page.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.NewHeader(gtx, theme)
		}),
		// Render the status panel.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.status.Draw(gtx, theme, statusMsgs, processing, done, total)
		}),
	)
}
