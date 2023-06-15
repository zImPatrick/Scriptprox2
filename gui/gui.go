package gui

import (
	"regexp"
	"scriptprox/gui/widgets"

	g "github.com/AllenDang/giu"
)

var (
	status           string
	serverCleanRegex *regexp.Regexp
	showServerlist   bool
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
				g.Checkbox("Serverliste anzeigen", &showServerlist),
				g.Condition(showServerlist, (func() g.Layout { //???????????
					if showServerlist {
						return widgets.ServerlistWidget(wnd, serverCleanRegex)
					} else {
						return g.Layout{}
					}
				})(), nil),
			),
			g.TabItem("Einstellungen").Layout(widgets.ConfigWidget()),
		),
	)
}

func RunGUI() {
	go widgets.RefreshServerlist()
	serverCleanRegex = regexp.MustCompile(`\^[0-9]`)
	wnd = g.NewMasterWindow("Scriptprox", 800, 400, g.MasterWindowFlags(g.WindowFlagsAlwaysAutoResize))
	wnd.Run(loop)
}
