package errorpage

import (
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/manager"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Page displays the interface when critical errors occur.
type Page struct {
	manager      *manager.Manager
	status       *components.Status
	SourceButton *widget.Clickable
}

func NewPage(mgr *manager.Manager, statusComponent *components.Status) *Page {
	return &Page{
		manager:      mgr,
		status:       statusComponent,
		SourceButton: &widget.Clickable{},
	}
}

func (p *Page) Draw(gtx layout.Context, theme *material.Theme) {
	// Retrieve the log history to display the error message to the user.
	statusMsgs, _, processing, _, done, total := p.manager.GetStatus()

	// Handle the back button to reset state and return to the start screen.
	if p.SourceButton.Clicked(gtx) {
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
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					// Back button to return to the main screen.
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, p.SourceButton, "Main menu")
						button.Inset = layout.UniformInset(unit.Dp(9))
						button.CornerRadius = unit.Dp(15)
						return button.Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.status.Draw(gtx, theme, statusMsgs, processing, done, total)
		}),
	)
}
