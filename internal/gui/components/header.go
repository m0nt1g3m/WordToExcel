package components

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"wordtoexcel/internal/gui/colors"
	"wordtoexcel/internal/gui/widgets"
)

func NewHeader(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	inset := layout.UniformInset(unit.Dp(10))
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, WeightSum: 1}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Width: unit.Dp(20), Height: unit.Dp(1)}.Layout(gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(unit.Dp(35))
				gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(35))
				return widgets.NewImage(gtx, "logo.png", layout.Center, 0)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Width: unit.Dp(10), Height: unit.Dp(1)}.Layout(gtx)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return widgets.NewLabel(
					gtx, theme, "Word to Excel", 22, colors.ColorTextMain, text.Start, font.Bold,
				)
			}),
		)
	},
	)
}
