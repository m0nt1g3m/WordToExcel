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
	welcomePage    *welcome.MainPage
	processingPage *processing.Page
	readyPage      *ready.Page
	errorPage      *errorpage.Page
	statusComp     *components.Status
	Manager        *manager.Manager
	Theme          *material.Theme
}

func NewWindow() *WordToExcel {
	theme := material.NewTheme()
	theme.Palette.ContrastBg = colors.ColorActionMain
	return &WordToExcel{
		Theme:   theme,
		Manager: manager.NewManager(),
	}
}

func (wte *WordToExcel) Init() {
	wte.Theme = material.NewTheme()
	wte.Theme.Palette.ContrastBg = colors.ColorActionMain
	wte.Theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	wte.Theme.Face = font.Typeface("Go")

	wte.Manager.Window.Option(
		app.Title("WordToExcel"),
		app.Size(unit.Dp(640), unit.Dp(480)),
		app.MaxSize(unit.Dp(1280), unit.Dp(720)),
		app.MinSize(unit.Dp(1280), unit.Dp(720)),
	)
	wte.initPages()
}

func (wte *WordToExcel) initPages() {
	wte.statusComp = components.NewStatus()
	wte.processingPage = processing.NewPage(wte.Manager, wte.statusComp)
	wte.welcomePage = welcome.NewMainPage(wte.Manager, wte.statusComp)
	wte.readyPage = ready.NewPage(wte.Manager, wte.statusComp)
	wte.errorPage = errorpage.NewPage(wte.Manager, wte.statusComp)
}

func (wte *WordToExcel) Draw() {
	for {
		evt := wte.Manager.Window.Event()

		switch e := evt.(type) {
		case app.DestroyEvent:
			wte.Manager.Cleanup()
			os.Exit(0)
		case app.FrameEvent:
			var ops op.Ops
			gtx := app.NewContext(&ops, e)
			paint.ColorOp{Color: colors.ColorBackground}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			switch wte.Manager.State() {
			case vars.StateInit, vars.StateFilesSelected:
				wte.welcomePage.Draw(gtx, wte.Theme)
			case vars.StateProcessing:
				wte.processingPage.Draw(gtx, wte.Theme)
			case vars.StateReady:
				wte.readyPage.Draw(gtx, wte.Theme)
			case vars.StateError:
				wte.errorPage.Draw(gtx, wte.Theme)
			}
			e.Frame(gtx.Ops)
		}
	}
}
