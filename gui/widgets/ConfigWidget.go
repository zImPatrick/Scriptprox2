package widgets

import (
	"image/color"
	"scriptprox/settings"
	"scriptprox/utils"

	g "github.com/AllenDang/giu"
)

var config *settings.Settings

func ConfigWidget() g.Layout {
	if config == nil {
		config = settings.GetSettings()
	}

	return g.Layout{
		g.Checkbox("Tokens nicht erstellen/HWID Ban Bypass", &config.HWIDBypass),
		g.Tooltip("Generiert sogenannte \"Tokens\" nicht. Diese werden oft von Servern benutzt, um deine HWID zu bannen."),
		g.Row(
			g.Label("Resource-Name"),
			g.InputText(&config.ResourceName),
		),
		g.Tooltip("So \"tarnt\" sich das Mod-Menu im Spiel. spawnmanager ist oft ein gutes \"Versteck\"."),
		g.Condition(config.ProfiModus, g.Layout{
			g.Separator(),
			g.Row(
				g.Label("GUID Override"),
				g.InputText(&config.GUIDOverride),
			),
			g.Row(
				g.Label("Ticket Override"),
				g.InputText(&config.TicketOverride),
			),
		}, g.Layout{}),
		g.Button("Speichern").OnClick(settings.SaveSettings),
		g.Style().SetColor(g.StyleColorText, color.Gray{Y: 80}).To(
			g.Label("Version: " + utils.Commit),
		),
	}
}
