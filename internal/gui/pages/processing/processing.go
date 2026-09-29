package processing

import (
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/manager"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Page displays the interface during the conversion process.
type Page struct {
	manager *manager.Manager
	status  *components.Status
}

// NewPage creates and returns a new processing Page.
func NewPage(mgr *manager.Manager, status *components.Status) *Page {
	return &Page{
		manager: mgr,
		status:  status,
	}
}

// Draw renders the page for the current frame.
func (p *Page) Draw(gtx layout.Context, theme *material.Theme) {
	msgs, _, processing, _, done, total := p.manager.GetStatus()

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.NewHeader(gtx, theme)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.status.Draw(gtx, theme, msgs, processing, done, total)
		}),
	)
}
