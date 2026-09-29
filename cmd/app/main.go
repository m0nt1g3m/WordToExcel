package main

import (
	"wordtoexcel/internal/gui/pages"

	"gioui.org/app"
)

var WordToExcel pages.WordToExcel

func main() {
	WordToExcel = *pages.NewWindow()

	WordToExcel.Init()

	go WordToExcel.Draw()

	app.Main()
}
