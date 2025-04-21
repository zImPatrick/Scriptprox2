package gui

import (
	_ "embed"
	"regexp"
	"scriptprox/gui/resources"
	"scriptprox/gui/widgets"

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
			g.TabItem("Servers").Layout(
				widgets.ManualInputWidget(&status, serverCleanRegex),
				g.Condition(status != "", g.Layout{
					g.Label("Status: " + status),
				}, nil),
				widgets.ServerlistWidget(wnd, serverCleanRegex),
			),
			g.TabItem("Settings").Layout(widgets.ConfigWidget()),
		),
	)
}

//go:embed resources/forkawesome-webfont.ttf
var IconFont []byte

func RunGUI() {
	serverCleanRegex = regexp.MustCompile(`\^[0-9]`)
	go widgets.RefreshServerlist()
	wnd = g.NewMasterWindow("Scriptprox", 800, 400, g.MasterWindowFlags(g.WindowFlagsAlwaysAutoResize))
	font := g.Context.FontAtlas.AddFontFromBytes("Icons", IconFont, 16)
	resources.IconFont = font
	wnd.Run(loop)
}
