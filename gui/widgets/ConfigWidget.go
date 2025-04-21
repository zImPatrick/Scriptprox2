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
		g.Checkbox("Do not create tokens", &config.HWIDBypass),
		g.Tooltip("Do not generate tokens, this is used by servers to ban by HWID"),

		g.Row(
			g.Label("Resource-Name"),
			g.InputText(&config.ResourceName),
		),
		g.Tooltip("The resource name your resource.rpf is using (spawnmanager is recommended)"),
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
		g.Button("Save").OnClick(settings.SaveSettings),
		g.Style().SetColor(g.StyleColorText, color.Gray{Y: 80}).To(
			g.Label("Version: " + utils.Commit),
		),
	}
}
