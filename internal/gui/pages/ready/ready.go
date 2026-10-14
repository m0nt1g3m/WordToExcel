package ready

import (
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/manager"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Page handles rendering the successful completion screen after conversion.
type Page struct {
	manager      *manager.Manager
	status       *components.Status
	SourceButton *widget.Clickable
	SaveButton   *widget.Clickable
}

func NewPage(mgr *manager.Manager, statusComponent *components.Status) *Page {
	return &Page{
		manager:      mgr,
		status:       statusComponent,
		SourceButton: &widget.Clickable{},
		SaveButton:   &widget.Clickable{},
	}
}

func (p *Page) Draw(gtx layout.Context, theme *material.Theme) {
	statusMsgs, _, processing, readyToSave, done, total := p.manager.GetStatus()

	// Handle the request to return to the initial screen.
	if p.SourceButton.Clicked(gtx) {
		p.manager.Reset()
	}
	// Handle the request to save results if they are ready to be exported.
	if readyToSave && p.SaveButton.Clicked(gtx) {
		p.manager.SaveResult()
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
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					// Back button to return to the main screen.
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, p.SourceButton, "Main menu")
						button.Inset = layout.UniformInset(unit.Dp(9))
						button.CornerRadius = unit.Dp(15)
						return button.Layout(gtx)
					}),
					// Save file button shown only until the save operation succeeds.
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !readyToSave {
							return layout.Dimensions{}
						}
						return layout.Inset{Left: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							button := material.Button(theme, p.SaveButton, "Save result")
							button.Inset = layout.UniformInset(unit.Dp(9))
							button.CornerRadius = unit.Dp(15)
							return button.Layout(gtx)
						})
					}),
				)
			})
		}),
		// Display the status panel to provide the user with a final conversion report.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.status.Draw(gtx, theme, statusMsgs, processing, done, total)
		}),
	)
}
