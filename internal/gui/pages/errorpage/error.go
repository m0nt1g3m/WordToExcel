package errorpage

import (
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/gui/widgets"
	"wordtoexcel/internal/manager"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Page displays the interface in case of critical errors.
type Page struct {
	manager      *manager.Manager
	status       *components.Status
	sourceButton *widget.Clickable
}

// NewPage creates and returns a new error Page.
func NewPage(mgr *manager.Manager, status *components.Status) *Page {
	return &Page{
		manager:      mgr,
		status:       status,
		sourceButton: &widget.Clickable{},
	}
}

// Draw renders the page for the current frame.
func (p *Page) Draw(gtx layout.Context, theme *material.Theme) {
	msgs, _, processing, _, done, total := p.manager.GetStatus()

	// "Home" button resets app to initial state.
	if p.sourceButton.Clicked(gtx) {
		p.manager.Reset()
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.NewHeader(gtx, theme)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return widgets.NewButton(gtx, theme, p.sourceButton, "Home")
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.status.Draw(gtx, theme, msgs, processing, done, total)
		}),
	)
}
