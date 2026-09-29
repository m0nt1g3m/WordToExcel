package pages

import (
	"os"
	"wordtoexcel/internal/gui/colors"
	"wordtoexcel/internal/gui/components"
	"wordtoexcel/internal/gui/pages/errorpage"
	"wordtoexcel/internal/gui/pages/processing"
	"wordtoexcel/internal/gui/pages/ready"
	"wordtoexcel/internal/gui/pages/welcome"
	"wordtoexcel/internal/manager"
	"wordtoexcel/internal/vars"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type WordToExcel struct {
	Manager *manager.Manager
	Theme   *material.Theme
	routes  map[vars.AppState]Page
}

func NewWindow() *WordToExcel {
	return &WordToExcel{
		Theme:   material.NewTheme(),
		Manager: manager.NewManager(),
	}
}

// Init configures the window and registers all pages.
func (wte *WordToExcel) Init() {
	wte.initTheme()
	wte.initWindow()
	wte.initPages()
}

func (wte *WordToExcel) initTheme() {
	wte.Theme = material.NewTheme()
	wte.Theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	wte.Theme.Face = font.Typeface("Go")
}

func (wte *WordToExcel) initWindow() {
	wte.Manager.Window.Option(
		app.Title("WordToExcel"),
		app.MinSize(unit.Dp(800), unit.Dp(600)),
		app.MaxSize(unit.Dp(1024), unit.Dp(768)),
		app.Size(unit.Dp(800), unit.Dp(600)),
	)
}

// initPages constructs all pages and maps them to application states.
func (wte *WordToExcel) initPages() {
	status := components.NewStatus()

	wte.routes = map[vars.AppState]Page{
		vars.StateInit:          welcome.NewPage(wte.Manager, status),
		vars.StateFilesSelected: welcome.NewPage(wte.Manager, status),
		vars.StateProcessing:    processing.NewPage(wte.Manager, status),
		vars.StateReady:         ready.NewPage(wte.Manager, status),
		vars.StateError:         errorpage.NewPage(wte.Manager, status),
	}
}

// Draw is the main event loop. It dispatches rendering to the active page.
func (wte *WordToExcel) Draw() {
	for {
		evt := wte.Manager.Window.Event()

		switch e := evt.(type) {
		case app.DestroyEvent:
			os.Exit(0)

		case app.FrameEvent:
			var ops op.Ops
			gtx := app.NewContext(&ops, e)

			// Fill background.
			paint.ColorOp{Color: colors.ColorBackground}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// Dispatch to the active page.
			if page, ok := wte.routes[wte.Manager.State()]; ok {
				page.Draw(gtx, wte.Theme)
			}

			e.Frame(gtx.Ops)
		}
	}
}
