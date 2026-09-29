package ready

import (
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/gui/widgets"
	"wordtoexcel/internal/manager"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Page is the success screen shown after all files are processed without critical errors.
type Page struct {
	manager      *manager.Manager
	status       *components.Status
	sourceButton *widget.Clickable
	saveButton   *widget.Clickable
}

// NewPage creates and returns a new ready Page.
func NewPage(mgr *manager.Manager, status *components.Status) *Page {
	return &Page{
		manager:      mgr,
		status:       status,
		sourceButton: &widget.Clickable{},
		saveButton:   &widget.Clickable{},
	}
}

// Draw renders the page for the current frame.
func (p *Page) Draw(gtx layout.Context, theme *material.Theme) {
	msgs, _, processing, readyToSave, done, total := p.manager.GetStatus()

	// "Home" resets to the initial state.
	if p.sourceButton.Clicked(gtx) {
		p.manager.Reset()
	}
	// "Save Result" is only active while the temp file is available.
	if readyToSave && p.saveButton.Clicked(gtx) {
		p.manager.SaveResult()
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.NewHeader(gtx, theme)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
		}),
		// Action buttons row.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return widgets.NewButton(gtx, theme, p.sourceButton, "Home")
					}),
					// "Save Result" appears only until the file is saved.
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !readyToSave {
							return layout.Dimensions{}
						}
						return layout.Inset{Left: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return widgets.NewButton(gtx, theme, p.saveButton, "Save Result")
						})
					}),
				)
			})
		}),
		// Status / log panel with the final conversion report.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.status.Draw(gtx, theme, msgs, processing, done, total)
		}),
	)
}
