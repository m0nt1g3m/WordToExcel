package components

import (
	"fmt"
	"image"
	"wordtoexcel/internal/gui/colors"
	"wordtoexcel/internal/gui/fonts"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const statusBoxHeight = unit.Dp(200)

// Status is a shared log + progress component used by all pages.
type Status struct {
	list layout.List
}

func NewStatus() *Status {
	return &Status{}
}

func (s *Status) Draw(gtx layout.Context, theme *material.Theme, messages []string, processing bool, done, total int) layout.Dimensions {
	if len(messages) == 0 {
		return layout.Dimensions{}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Progress counter (hidden when total == 0).
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if total == 0 {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(8), Right: unit.Dp(48), Left: unit.Dp(48)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(theme, fonts.SizeCaption, fmt.Sprintf("Processed: %d of %d", done, total))
				label.Color = colors.ColorTextSub
				return label.Layout(gtx)
			})
		}),
		// Log panel with border.
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(18), Right: unit.Dp(48), Bottom: unit.Dp(18), Left: unit.Dp(48)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				width := gtx.Constraints.Max.X
				height := gtx.Metric.Dp(statusBoxHeight)
				if gtx.Constraints.Max.Y < height {
					height = gtx.Constraints.Max.Y
				}
				gtx.Constraints.Min = image.Pt(width, height)
				gtx.Constraints.Max = image.Pt(width, height)

				return layout.Stack{}.Layout(gtx,
					// Border + background.
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						outer := clip.RRect{Rect: image.Rectangle{Max: gtx.Constraints.Max}, SE: 12, SW: 12, NE: 12, NW: 12}
						paint.FillShape(gtx.Ops, colors.ColorActionMain, outer.Op(gtx.Ops))
						inner := clip.RRect{Rect: image.Rectangle{Min: image.Pt(1, 1), Max: image.Pt(width-1, height-1)}, SE: 11, SW: 11, NE: 11, NW: 11}
						paint.FillShape(gtx.Ops, colors.ColorBackgroundLight, inner.Op(gtx.Ops))
						return layout.Dimensions{Size: gtx.Constraints.Max}
					}),
					// Content.
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(14), Right: unit.Dp(16), Bottom: unit.Dp(10), Left: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								// Section title.
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Label(theme, fonts.SizeCaption, "Processing Status")
									label.Color = colors.ColorTextSub
									return label.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Spacer{Height: unit.Dp(8)}.Layout(gtx)
								}),
								// Scrollable log messages.
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									s.list.Axis = layout.Vertical
									s.list.ScrollToEnd = processing
									return s.list.Layout(gtx, len(messages), func(gtx layout.Context, i int) layout.Dimensions {
										label := material.Label(theme, fonts.SizeCaption, messages[i])
										label.Color = colors.ColorTextMain
										return label.Layout(gtx)
									})
								}),
								// "Processing continues..." spinner label.
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if !processing {
										return layout.Dimensions{}
									}
									label := material.Label(theme, fonts.SizeSmall, "Processing continues...")
									label.Color = colors.ColorTextSub
									return label.Layout(gtx)
								}),
							)
						})
					}),
				)
			})
		}),
	)
}

