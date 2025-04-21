package widgets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"scriptprox/gui/resources"
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
	serverIdOrUrl      string
	server             ServerInfo
	manualServerIp     string
	profiModusDropdown bool
)

func FindServer(status *string, idOrUrl string) {
	*status = ""
	serverIdOrUrl = strings.TrimPrefix(idOrUrl, "cfx.re/join/")

	resp, err := http.Get("https://servers-frontend.fivem.net/api/servers/single/" + serverIdOrUrl)
	if err != nil {
		*status = "Couldn't retrieve server (" + err.Error() + ")"
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		*status = "Couldn't retrieve server (unable to read body, " + err.Error() + ")"
		return
	}
	err = json.Unmarshal(respBody, &server)
	if err != nil {
		*status = "Couldn't retrieve server (failed parsing JSON, " + err.Error() + ")"
		return
	}
	*status = ""
}

func connectToServer(status *string) {
	connectEndpoint := server.Data.ConnectEndpoints[0]
	if !strings.HasPrefix(connectEndpoint, "http") {
		connectEndpoint = "https://" + connectEndpoint
	}
	httpScriptprox.ChangeEndpoint(connectEndpoint)
	exec.Command("rundll32", "url.dll,FileProtocolHandler", "fivem://connect/localhost:30120").Run()
	*status = "Connecting"
}

func ManualInputWidget(status *string, serverCleanRegex *regexp.Regexp) g.Layout {
	if config == nil {
		config = settings.GetSettings()
	}
	return g.Layout{
		g.Row(
			g.InputText(&serverIdOrUrl).Hint("cfx.re/join/zq4ayd").Size(200),
			resources.WithIconFont(
				g.Row(
					g.Button("\uf002").OnClick(func() { go FindServer(status, serverIdOrUrl) }),
					g.Button("\uf0ea").OnClick(func() { go FindServer(status, g.Context.GetPlatform().GetClipboard()) }),
					g.Condition(config.ProfiModus, g.Layout{
						g.Button("\uf0dd").OnClick(func() { profiModusDropdown = !profiModusDropdown }),
					}, g.Layout{
						g.Dummy(0, 0),
					}),
				),
			),
		),
		g.Condition(profiModusDropdown, g.Layout{
			g.Row(
				g.InputText(&manualServerIp).Hint("http://127.0.0.1:30120"),
				g.Button("Connect").OnClick(func() {
					err := httpScriptprox.ChangeEndpoint(manualServerIp)
					if err != nil {
						*status = err.Error()
					} else {
						exec.Command("rundll32", "url.dll,FileProtocolHandler", "fivem://connect/localhost:30120").Run()
						*status = "Connecting"
					}
				}),
			),
		}, g.Layout{}),
		g.Condition(server.EndPoint != "", g.Layout{
			g.Row(
				// this... doesn't look great, but it works
				g.Condition(server.Data.IconVersion != 0, g.Layout{
					g.ImageWithURL("https://servers-live.fivem.net/servers/icon/" + server.EndPoint + "/" + fmt.Sprint(server.Data.IconVersion) + ".png").LayoutForLoading(
						g.Dummy(96, 96),
					).LayoutForFailure(
						g.Dummy(96, 96),
					),
				}, g.Layout{g.Dummy(96, 96)}),
				g.Column(
					g.Label(serverCleanRegex.ReplaceAllString(server.Data.Hostname, "")),
					g.Row(
						g.Label("\uf007").Font(resources.IconFont),
						g.Label(fmt.Sprint(server.Data.Clients)+"/"+fmt.Sprint(server.Data.SvMaxclients)),
					),
					g.Tooltip("Amount of players currently online"),
					g.Row(
						g.Button("Connect").OnClick(func() { connectToServer(status) }),
						resources.WithIconFont(
							g.Row(
								g.Button("\uf0c5").OnClick(func() { g.Context.GetPlatform().SetClipboard(server.EndPoint) }),
								g.Button("\uf0dd").OnClick(func() { // Verstecken
									server = ServerInfo{}
								}),
							),
						),
					),
				),
			),
		}, nil),
	}
}
