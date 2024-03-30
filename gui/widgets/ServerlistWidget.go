package widgets

import (
	"fmt"
	"regexp"
	"scriptprox/gui/resources"
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
	rows := make([]*g.TableRowWidget, len(serversInServerlist))

	amount := 0
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

		// Workaround für Bug in Golang: https://go.dev/blog/loopvar-preview
		// kb gerade auf 1.21 zu migraten
		serverInListIndex := i
		rows[amount] = g.TableRow(
			g.Condition(serversInServerlist[i].Data.IconVersion != 0, g.Layout{
				g.ImageWithURL(
					"https://servers-live.fivem.net/servers/icon/"+serversInServerlist[i].EndPoint+"/"+fmt.Sprint(serversInServerlist[i].Data.IconVersion)+".png",
				).Size(48, 48).LayoutForFailure(g.Dummy(48, 48)).LayoutForLoading(g.Dummy(48, 48)),
			}, g.Layout{g.Align(g.AlignCenter).To(g.Dummy(48, 48))}),
			g.Label(unfilteredName),
			g.Custom(func() {
				w, _ := wnd.GetSize()
				g.Selectable(lang).Size(float32(w), 48).Flags(g.SelectableFlagsSpanAllColumns).OnClick(func() {
					unnessecaryString := ""
					go ServerSuchen(&unnessecaryString, serversInServerlist[serverInListIndex].EndPoint)
				}).OnDClick(func() {
					go func() {
						unnessecaryString := ""
						ServerSuchen(&unnessecaryString, serversInServerlist[serverInListIndex].EndPoint)
						verbinden(&unnessecaryString)
					}()
				}).Build()
			}),
			g.Label(strconv.Itoa(int(serversInServerlist[i].Data.Clients))),
		).MinHeight(48)
		amount = amount + 1
	}
	return rows[:amount]
}

var rows []*g.TableRowWidget

func ServerlistWidget(wnd *g.MasterWindow, serverCleanRegex *regexp.Regexp) g.Layout {
	if !isServerlistLoading && rows == nil {
		fmt.Println("updating server list")

		// warum muss ich rows hier nochmal setzen??
		// ich mach das doch in der methode?
		// blöder computer
		rows = buildServerlistRows(wnd, serverCleanRegex)
	}

	return g.Layout{
		g.Row(
			g.InputText(&Filters.servername).OnChange(func() {
				rows = buildServerlistRows(wnd, serverCleanRegex)
			}).Hint("Nach Server suchen"),
			g.InputText(&Filters.land).OnChange(func() {
				rows = buildServerlistRows(wnd, serverCleanRegex)
			}).Hint("Land eingeben (z.B. de, fr)").Size(220),
			resources.WithIconFont(
				g.Button("\uf021").OnClick(func() { go RefreshServerlist() }),
			),
		),
		g.Table().FastMode(true).Freeze(0, 1).Columns(
			g.TableColumn("").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(48),
			g.TableColumn("Server"),
			g.TableColumn("Land").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(32),
			g.TableColumn("Spieler").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(48),
		).Rows(rows...),
	}
}

var isServerlistLoading bool

func RefreshServerlist() {
	if isServerlistLoading {
		// wtf? wir können jetzt nicht einfach reloaden während schon geloaded wird
		return
	}
	isServerlistLoading = true
	serversInServerlist = []*serverlistProtos.Server{}
	rows = nil
	serverlist.RequestServerlist(func(s *serverlistProtos.Server) {
		serversInServerlist = append(serversInServerlist, s)
	}, func() {
		// fertig
		sort.SliceStable(serversInServerlist, func(i, j int) bool {
			return serversInServerlist[i].Data.Clients > serversInServerlist[j].Data.Clients
		})
		isServerlistLoading = false
		rows = nil
	})

}
