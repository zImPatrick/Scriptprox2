package widgets

import (
	"scriptprox/settings"

	g "github.com/AllenDang/giu"
)

var config *settings.Settings

func ConfigWidget() g.Layout {
	if config == nil {
		config = settings.GetSettings()
	}

	return g.Layout{
		g.Checkbox("Tokens nicht erstellen/HWID Bypass", &config.HWIDBypass),
		g.Row(
			g.Label("Resource-Name"),
			g.InputText(&config.ResourceName),
		),
		g.Button("Speichern").OnClick(settings.SaveSettings),
	}
}
