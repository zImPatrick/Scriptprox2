package widgets

import (
	"fmt"
	"net/url"
	"os/exec"
	"scriptprox/gui/resources"
	"scriptprox/settings"
	"scriptprox/updater"

	g "github.com/AllenDang/giu"
)

func GetStringForUpdate(update updater.UpdateInfo) string {
	switch update.State {
	case updater.AVAILABLE:
		return "Ein Update ist verfügbar!"
	case updater.UPDATING:
		return "Update wird durchgeführt..."
	case updater.DONE:
		if len(update.EncounteredErrors) != 0 {
			return "Das Update wurde mit Fehlern durchgeführt. Mehr Details befinden sich unten."
		}

		return "Das Update wurde erfolgreich durchgeführt."
	default:
		return fmt.Sprintf("Update keinen String für State %d, bitte melden", update.State)
	}
}

func UpdateWidget() g.Layout {
	if config.FuckUpdates {
		return g.Layout{
			g.Label("Du hast die Einstellung 'FuckUpdates' aktiviert. Willst du Sie deaktivieren?"),
			g.Button("Ja").OnClick(func() {
				config.FuckUpdates = false
				settings.SaveSettings()

				go updater.CheckForUpdates()
			}),
		}
	}
	if updater.LastUpdate == nil { // sollten vielleicht auch error ausgeben
		return g.Layout{
			g.Label("Suche nach Updates..."),
		}
	}
	if updater.LastUpdate.State == updater.HOST_IS_FUCKED {
		return g.Layout{
			g.Label("Alle Updater-Hosts sind tot, schreib vllt mal Patrick"),
		}
	}
	if updater.LastUpdate.State == updater.NONEXISTANT {
		return g.Layout{
			g.Label("Kein Update ist derzeit verfügbar."),
		}
	} else {
		return g.Layout{
			g.Row(
				g.Button("Updaten").OnClick(func() { go updater.InstallUpdate(updater.LastUpdate) }).Disabled(updater.LastUpdate == nil ||
					updater.LastUpdate.State != updater.AVAILABLE),
				g.Label(GetStringForUpdate(*updater.LastUpdate)),
			),
			g.Table().Columns(
				g.TableColumn("Datei"),
				g.TableColumn("Fehler"),
				g.TableColumn("Behebung"),
			).Rows(
				func() []*g.TableRowWidget {
					list := make([]*g.TableRowWidget, 0)
					for _, v := range updater.LastUpdate.RequiredFiles {
						widgets := []g.Widget{
							g.Row(
								resources.WithIconFont(
									g.Label("\uf15b"),
								),
								g.Label(v),
							),
						}
						if err, ok := updater.LastUpdate.EncounteredErrors[v]; ok {
							widgets = append(widgets, g.Label(err.Error()), g.Button("Manuell runterladen").OnClick(func() {
								url, _ := url.JoinPath(updater.UPDATER_HOST, v)
								exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
							}))
						}

						list = append(list, g.TableRow(widgets...))
					}

					return list
				}()...,
			),
		}
	}
}
