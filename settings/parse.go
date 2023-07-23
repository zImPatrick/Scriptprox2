package settings

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Settings struct {
	HWIDBypass     bool
	ResourceName   string
	GUIDOverride   string
	TicketOverride string
	ProfiModus	bool
}

var settings *Settings

func GetSettings() *Settings {
	if settings == nil {
		parseSettings()
	}

	return settings
}

func SaveSettings() {
	fileHandle, err := os.OpenFile("settings.toml", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Fehler beim Öffnen der Settings-Datei! " + err.Error())
	}
	err = toml.NewEncoder(fileHandle).Encode(settings)
	if err != nil {
		fmt.Println("Fehler beim Speichern der Settings-Datei! " + err.Error())
	}
}

func parseSettings() {
	// existiert die config?
	_, err := os.Stat("settings.toml")
	if err != nil {
		fmt.Println("Settings-Datei existiert nicht.")
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
			fmt.Println("Fehler beim Lesen der Settings-Datei! " + err.Error())
			os.Exit(1)
		}
	}
}
