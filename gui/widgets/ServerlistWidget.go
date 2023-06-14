package widgets

import (
	"fmt"
	"regexp"
	"scriptprox/gui/serverlist"
	serverlistProtos "scriptprox/gui/serverlist/protos"
	"sort"
	"strconv"
	"strings"

	g "github.com/AllenDang/giu"
)

var Filters struct {
	servername string
	land       string
}

var serversInServerlist []*serverlistProtos.Server

func buildServerlistRows(wnd *g.MasterWindow, serverCleanRegex *regexp.Regexp) []*g.TableRowWidget {
	rows := []*g.TableRowWidget{}
	// rows := make([]*g.TableRowWidget, len(serversInServerlist))
	for i := range serversInServerlist {
		unfilteredName := serverCleanRegex.ReplaceAllString(serversInServerlist[i].Data.Hostname, "")
		if Filters.servername != "" && !strings.Contains(strings.ToLower(unfilteredName), strings.ToLower(Filters.servername)) {
			continue
		}
		var lang string
		for varIndex := range serversInServerlist[i].Data.Vars {
			if serversInServerlist[i].Data.Vars[varIndex].Key == "locale" {
				lang = serversInServerlist[i].Data.Vars[varIndex].Value
				break
			}
		}
		if Filters.land != "" && !strings.Contains(lang, Filters.land) {
			continue
		}
		test := i
		rows = append(rows, g.TableRow(
			g.Condition(serversInServerlist[i].Data.IconVersion != 0, g.Layout{
				g.ImageWithURL("https://servers-live.fivem.net/servers/icon/"+serversInServerlist[i].EndPoint+"/"+fmt.Sprint(serversInServerlist[i].Data.IconVersion)+".png").Size(48, 48),
			}, g.Layout{g.Align(g.AlignCenter).To(g.Label("N/A"))}),
			g.Label(unfilteredName),
			g.Custom(func() {
				w, _ := wnd.GetSize()
				g.Selectable(lang).Size(float32(w), 48).Flags(g.SelectableFlagsSpanAllColumns).OnClick(func() {
					unnessecaryString := ""
					go ServerSuchen(&unnessecaryString, serversInServerlist[test].EndPoint)
				}).Build()
			}),
			g.Label(strconv.Itoa(int(serversInServerlist[i].Data.Clients))),
		).MinHeight(48))
	}

	return rows
}

func ServerlistWidget(wnd *g.MasterWindow, serverCleanRegex *regexp.Regexp) g.Layout {
	return g.Layout{
		g.Row(g.InputText(&Filters.servername).Hint("Nach Server suchen"), g.InputText(&Filters.land).Hint("Land eingeben (z.B. de, fr)")),
		g.Table().FastMode(true).Freeze(0, 1).Columns(
			g.TableColumn("").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(48),
			g.TableColumn("Server"),
			g.TableColumn("Land").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(32),
			g.TableColumn("Spieler").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(48),
		).Rows(buildServerlistRows(wnd, serverCleanRegex)...),
	}
}

func RefreshServerlist() {
	serversInServerlist = []*serverlistProtos.Server{}
	serverlist.RequestServerlist(func(s *serverlistProtos.Server) {
		if s.Data.Clients > 0 {
			serversInServerlist = append(serversInServerlist, s)
		}
	}, func() {
		// fertig
		sort.SliceStable(serversInServerlist, func(i, j int) bool {
			return serversInServerlist[i].Data.Clients > serversInServerlist[j].Data.Clients
		})
	})
}
