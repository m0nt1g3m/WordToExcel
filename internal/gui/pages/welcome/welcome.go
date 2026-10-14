package welcome

import (
	"fmt"
	"wordtoexcel/internal/gui/colors"
	"wordtoexcel/internal/gui/components"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (w *MainPage) Draw(gtx layout.Context, theme *material.Theme) {
	// Handle user input events.
	if w.SourceButton.Clicked(gtx) {
		go w.manager.SelectFiles()
	}
	if w.FolderButton.Clicked(gtx) {
		go w.manager.SelectFolder()
	}

	if w.ConvertButton.Clicked(gtx) {
		w.manager.ConvertFiles()
	}

	// Get the current state for rendering.
	selectedCount := w.manager.SelectedFilesCount()
	statusMsgs, _, processing, _, done, total := w.manager.GetStatus()

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Render the application header.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.NewHeader(gtx, theme)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						button := material.Button(theme, w.SourceButton, "Select files")
						button.Inset = layout.UniformInset(unit.Dp(9))
						button.CornerRadius = unit.Dp(15)
						return button.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							button := material.Button(theme, w.FolderButton, "Select folder")
							button.Inset = layout.UniformInset(unit.Dp(9))
							button.CornerRadius = unit.Dp(15)
							return button.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if selectedCount == 0 {
							return layout.Dimensions{}
						}
						return layout.Inset{Left: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							button := material.Button(theme, w.ConvertButton, "Convert")
							button.Inset = layout.UniformInset(unit.Dp(9))
							button.CornerRadius = unit.Dp(15)
							return button.Layout(gtx)
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if selectedCount == 0 {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(6), Left: unit.Dp(48)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(theme, unit.Sp(16), fmt.Sprintf("Files selected: %d", selectedCount))
				label.Color = colors.ColorTextSub
				return label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.status.Draw(gtx, theme, statusMsgs, processing, done, total)
		}),
	)
}
