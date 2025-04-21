package settings

import (
	"fmt"
	"os"
	"scriptprox/utils"

	"github.com/BurntSushi/toml"
)

type Settings struct {
	HWIDBypass     bool
	ResourceName   string
	GUIDOverride   string
	TicketOverride string
	ProfiModus     bool
}

var settings *Settings

func GetSettings() *Settings {
	if settings == nil {
		parseSettings()
	}

	return settings
}

func SaveSettings() {
	fileHandle, err := os.OpenFile("settings.toml", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Couldn't open settings file: " + err.Error())
	}
	err = toml.NewEncoder(fileHandle).Encode(settings)
	if err != nil {
		fmt.Println("Couldn't save settings file: " + err.Error())
	}
}

func parseSettings() {
	_, err := os.Stat("settings.toml")
	if err != nil {
		f, _ := os.Create("settings.toml")
		f.Close()
		settings = &Settings{
			HWIDBypass:   false,
			ResourceName: "spawnmanager",
		}
		SaveSettings()
	} else {
		// config existiert; parsen!
		_, err := toml.DecodeFile("settings.toml", &settings)
		if err != nil {
			utils.ThrowErrorAndQuit("Failed reading the settings file. Please delete it and try again.", err)
			os.Exit(1)
		}
	}
}
