package welcome

import (
	"fmt"
	"wordtoexcel/internal/gui/colors"
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/gui/fonts"
	"wordtoexcel/internal/gui/widgets"
	"wordtoexcel/internal/manager"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Page is the welcome / file-selection screen.
type Page struct {
	sourceButton  *widget.Clickable
	convertButton *widget.Clickable
	manager       *manager.Manager
	status        *components.Status
}

// NewPage creates and returns a new welcome Page.
func NewPage(mgr *manager.Manager, status *components.Status) *Page {
	return &Page{
		sourceButton:  &widget.Clickable{},
		convertButton: &widget.Clickable{},
		manager:       mgr,
		status:        status,
	}
}

// Draw renders the page for the current frame.
func (p *Page) Draw(gtx layout.Context, theme *material.Theme) {
	// Handle user input.
	if p.sourceButton.Clicked(gtx) {
		go p.manager.SelectFiles()
	}
	if p.convertButton.Clicked(gtx) {
		p.manager.ConvertFiles()
	}

	// Snapshot state for this frame.
	selectedCount := p.manager.SelectedFilesCount()
	msgs, _, processing, _, done, total := p.manager.GetStatus()

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Header bar.
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
					// "Select Files" is always visible.
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return widgets.NewButton(gtx, theme, p.sourceButton, "Select Files")
					}),
					// "Convert" appears only after files are selected.
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if selectedCount == 0 {
							return layout.Dimensions{}
						}
						return layout.Inset{Left: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return widgets.NewButton(gtx, theme, p.convertButton, "Convert")
						})
					}),
				)
			})
		}),
		// Selected-files count label.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if selectedCount == 0 {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(6), Left: unit.Dp(48)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(theme, fonts.SizeInfo, fmt.Sprintf("Selected files: %d", selectedCount))
				label.Color = colors.ColorTextSub
				return label.Layout(gtx)
			})
		}),
		// Status / log panel.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.status.Draw(gtx, theme, msgs, processing, done, total)
		}),
	)
}
