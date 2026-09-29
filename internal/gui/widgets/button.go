package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ButtonStyle defines common visual parameters for all app buttons.
const (
	buttonInset        = unit.Dp(9)
	buttonCornerRadius = unit.Dp(15)
)

// NewButton renders a styled primary button with the given label.
func NewButton(gtx layout.Context, theme *material.Theme, clickable *widget.Clickable, label string) layout.Dimensions {
	btn := material.Button(theme, clickable, label)
	btn.Inset = layout.UniformInset(buttonInset)
	btn.CornerRadius = buttonCornerRadius
	return btn.Layout(gtx)
}
