package gui

import (
	_ "embed"
	"image"
	"regexp"
	"scriptprox/gui/resources"
	"scriptprox/gui/widgets"
	"scriptprox/updater"

	g "github.com/AllenDang/giu"
)

var (
	status           string
	serverCleanRegex *regexp.Regexp
)

var wnd *g.MasterWindow

func loop() {
	g.SingleWindow().Layout(
		g.TabBar().TabItems(
			g.TabItem("Allgemein").Layout(
				widgets.ManualInputWidget(&status, serverCleanRegex),
				g.Condition(status != "", g.Layout{
					g.Label("Status: " + status),
				}, nil),
				widgets.ServerlistWidget(wnd, serverCleanRegex),
			),
			g.TabItem("Einstellungen").Layout(widgets.ConfigWidget()),
			g.TabItem("Updater").Layout(widgets.UpdateWidget()),
		),
		g.Custom(func() {
			if updater.LastUpdate != nil && updater.LastUpdate.State == updater.AVAILABLE {
				w, _ := wnd.GetSize()
				point := g.GetCursorPos()
				text := "Ein Update ist verfügbar!"
				textW, _ := g.CalcTextSize(text)
				g.SetCursorPos(image.Pt(w-int(textW)-2, 8))
				g.Label("Ein Update ist verfügbar!").Wrapped(true).Build()
				g.SetCursorPos(point)
			}
		}),
	)
}

//go:embed resources/forkawesome-webfont.ttf
var IconFont []byte

func RunGUI() {
	go widgets.RefreshServerlist()
	serverCleanRegex = regexp.MustCompile(`\^[0-9]`)
	wnd = g.NewMasterWindow("Scriptprox", 800, 400, g.MasterWindowFlags(g.WindowFlagsAlwaysAutoResize))
	font := g.Context.FontAtlas.AddFontFromBytes("Icons", IconFont, 16)
	resources.IconFont = font
	wnd.SetCloseCallback(func() bool {
		// #11, Nutzer verweigern das Programm zu schließen während ein Update läuft
		// Hat schon mal installationen corrupted lol
		if updater.LastUpdate != nil {
			return updater.LastUpdate.State != updater.UPDATING
		}

		return true
	})
	wnd.Run(loop)
}
