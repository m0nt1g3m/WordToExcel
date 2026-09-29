package widgets

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	txt "gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func NewLabel(gtx layout.Context, theme *material.Theme, text string,
	size int, color color.NRGBA, alignment txt.Alignment, weight font.Weight) layout.Dimensions {
	label := material.Label(theme, unit.Sp(size), text)
	label.Color = color
	label.Font.Weight = weight
	label.Alignment = alignment
	return label.Layout(gtx)
}
