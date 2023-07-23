package widgets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	httpScriptprox "scriptprox/proxies/http"
	"scriptprox/settings"

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
	serverIdOrUrl string
	server        ServerInfo // zu faul ein interface zu erstellen
	manualServerIp string
)

func ServerSuchen(status *string, idOrUrl string) {
	*status = ""
	serverIdOrUrl = idOrUrl
	serverIdOrUrl = strings.TrimPrefix(serverIdOrUrl, "cfx.re/join/")

	resp, err := http.Get("https://servers-frontend.fivem.net/api/servers/single/" + serverIdOrUrl)
	if err != nil {
		*status = "Fehler 1 beim Abrufen des Servers (" + err.Error() + ")"
		return
	}
	respBody, err := io.ReadAll(resp.Body) // todo: error handling
	if err != nil {
		*status = "Fehler 2 beim Abrufen des Servers (konnte Antwort nicht lesen, " + err.Error() + ")"
		return
	}
	err = json.Unmarshal(respBody, &server)
	if err != nil {
		*status = "Fehler 3 beim Abrufen des Servers (ist kein JSON, " + err.Error() + ")"
		return
	}
	*status = ""
}

func verbinden(status *string) {
	connectEndpoint := server.Data.ConnectEndpoints[0]
	if !strings.HasPrefix(connectEndpoint, "http") {
		connectEndpoint = "https://" + connectEndpoint
	}
	httpScriptprox.ChangeEndpoint(connectEndpoint)
	exec.Command("rundll32", "url.dll,FileProtocolHandler", "fivem://connect/localhost:30120").Run()
	*status = "FiveM wurde gestartet."
}

func ManualInputWidget(status *string, serverCleanRegex *regexp.Regexp) g.Layout {
	if config == nil {
		config = settings.GetSettings()
	}
	return g.Layout{
		g.Row(
			g.Label("Server eingeben"),
			g.InputText(&serverIdOrUrl).Hint("cfx.re/join/zq4ayd"),
			g.Button("Server suchen").OnClick(func() { go ServerSuchen(status, serverIdOrUrl) }),
		),
		g.Row(
			g.Condition(config.ProfiModus, g.Layout{
				g.InputText(&manualServerIp).Hint("https://zimpatrick.gq"),
				g.Button("Einstellen").OnClick(func() {
					err := httpScriptprox.ChangeEndpoint(manualServerIp)
					if err != nil {
						*status = err.Error()
					} else {
						exec.Command("rundll32", "url.dll,FileProtocolHandler", "fivem://connect/localhost:30120").Run()
						*status = "FiveM wurde gestartet."
					}
				}),
			}, g.Layout{}),
		),
		g.Condition(server.EndPoint != "", g.Layout{
			g.Row(
				// mir gefällt das hier nicht.
				g.Condition(server.Data.IconVersion != 0, g.Layout{
					g.ImageWithURL("https://servers-live.fivem.net/servers/icon/" + server.EndPoint + "/" + fmt.Sprint(server.Data.IconVersion) + ".png").LayoutForLoading(
						g.Label("N/A"),
					).LayoutForFailure(
						g.Label("N/A"),
					),
				}, g.Layout{g.Label("N/A")}),
				g.Column(
					g.Label(serverCleanRegex.ReplaceAllString(server.Data.Hostname, "")),
					g.Label(fmt.Sprint(server.Data.Clients)+"/"+fmt.Sprint(server.Data.SvMaxclients)+" Spieler"),
					g.Button("Verbinden").OnClick(func() { verbinden(status) }),
				),
			),
		}, nil),
	}
}
