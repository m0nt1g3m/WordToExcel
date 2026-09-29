package pages

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Page is the interface every application screen must implement.
// To add a new page:
//  1. Create a package under internal/gui/pages/<yourpage>/
//  2. Define a struct with manager and status fields
//  3. Implement Draw(gtx, theme)
//  4. Register it in wordtoexcel.go: initPages() + route it in Draw()
type Page interface {
	Draw(gtx layout.Context, theme *material.Theme)
}
