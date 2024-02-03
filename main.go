package main

import (
	"fmt"
	"scriptprox/gui"
	"scriptprox/proxies/http"
	"scriptprox/proxies/udp"
	"scriptprox/settings"
	"scriptprox/updater"
	"sync"
)

func main() {
	go DoLicenseCheck()
	fmt.Println("Starting....")
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)
	go http.InitProxy(&waitgroup)
	go udp.InitProxy(&waitgroup)

	if !settings.GetSettings().FuckUpdates {
		go updater.CheckForUpdates()
	}

	gui.RunGUI()
}
