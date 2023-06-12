package gui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"scriptprox/gui/serverlist"
	serverlistProtos "scriptprox/gui/serverlist/protos"
	httpScriptprox "scriptprox/proxies/http"

	g "github.com/AllenDang/giu"
)

type ServerData struct {
	ConnectEndpoints []string
	Hostname         string
	Clients          int
	Vars             map[string]string
	IconVersion      int
	SvMaxclients     int
}

type ServerInfo struct {
	EndPoint string
	Data     ServerData
}

var (
	serverIdOrUrl    string
	status           string
	server           ServerInfo // zu faul ein interface zu erstellen
	serverCleanRegex *regexp.Regexp
)

var serversInServerlist []*serverlistProtos.Server

func serverSuchen() {
	status = ""
	serverIdOrUrl = strings.TrimPrefix(serverIdOrUrl, "cfx.re/join/")

	resp, err := http.Get("https://servers-frontend.fivem.net/api/servers/single/" + serverIdOrUrl)
	if err != nil {
		status = "Fehler 1 beim Abrufen des Servers (" + err.Error() + ")"
		return
	}
	respBody, err := io.ReadAll(resp.Body) // todo: error handling
	if err != nil {
		status = "Fehler 2 beim Abrufen des Servers (konnte Antwort nicht lesen, " + err.Error() + ")"
		return
	}
	err = json.Unmarshal(respBody, &server)
	if err != nil {
		status = "Fehler 3 beim Abrufen des Servers (ist kein JSON, " + err.Error() + ")"
		return
	}
	status = ""
}

func verbinden() {
	connectEndpoint := server.Data.ConnectEndpoints[0]
	if !strings.HasPrefix(connectEndpoint, "http") {
		connectEndpoint = "https://" + connectEndpoint
	}
	httpScriptprox.ChangeEndpoint(connectEndpoint)
	exec.Command("rundll32", "url.dll,FileProtocolHandler", "fivem://connect/localhost:30120").Run()
	status = "FiveM sollte nun gestartet werden. \nNutze:\n> connect localhost:30120\n in der FiveM-Konsole falls es nicht geht."
}

var Filters struct {
	servername string
	land       string
}

func buildServerlistRows() []*g.TableRowWidget {
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
			g.ImageWithURL("https://servers-live.fivem.net/servers/icon/"+serversInServerlist[i].EndPoint+"/"+fmt.Sprint(serversInServerlist[i].Data.IconVersion)+".png").Size(48, 48),
			g.Label(unfilteredName),
			g.Selectable(lang).Flags(g.SelectableFlagsSpanAllColumns).OnClick(func() {
				serverIdOrUrl = serversInServerlist[test].EndPoint
				serverSuchen()
			}),
			g.Label(strconv.Itoa(int(serversInServerlist[i].Data.Clients))),
		).MinHeight(48))
	}

	// rows[0].BgColor(&(color.RGBA{200, 100, 100, 255}))

	return rows
}

func loop() {
	g.SingleWindow().Layout(
		g.Row(
			g.Label("Server-URL eingeben"),
			g.InputText(&serverIdOrUrl),
			g.Button("Server suchen").OnClick(serverSuchen),
		),
		g.Condition(server.EndPoint != "", g.Layout{
			g.Row(
				g.Condition(server.Data.IconVersion != 0, g.Layout{
					g.ImageWithURL("https://servers-live.fivem.net/servers/icon/" + server.EndPoint + "/" + fmt.Sprint(server.Data.IconVersion) + ".png"),
				}, nil),
				g.Column(
					g.Label(serverCleanRegex.ReplaceAllString(server.Data.Hostname, "")),
					g.Label(fmt.Sprint(server.Data.Clients)+"/"+fmt.Sprint(server.Data.SvMaxclients)+" Spieler"),
					g.Button("Verbinden").OnClick(verbinden),
				),
			),
		}, nil),
		g.Condition(status != "", g.Layout{
			g.Label("Status: " + status),
		}, nil),
		g.Row(g.InputText(&Filters.servername).Hint("Nach Server suchen"), g.InputText(&Filters.land).Hint("Land eingeben (z.B. de, fr)")),
		g.Table().FastMode(true).Columns(
			g.TableColumn("").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(48),
			g.TableColumn("Server"),
			g.TableColumn("Land").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(32),
			g.TableColumn("Spieler").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(48),
		).Rows(buildServerlistRows()...),
	)
}

func RefreshServerlist() {
	serversInServerlist = []*serverlistProtos.Server{}
	serverlist.RequestServerlist(func(s *serverlistProtos.Server) {
		if s.Data.Clients > 0 {
			serversInServerlist = append(serversInServerlist, s)
		}
	})
}

func RunGUI() {
	go RefreshServerlist()
	serverCleanRegex = regexp.MustCompile(`\^[0-9]`)
	wnd := g.NewMasterWindow("Scriptprox", 800, 400, g.MasterWindowFlags(g.WindowFlagsAlwaysAutoResize))
	wnd.Run(loop)
}
