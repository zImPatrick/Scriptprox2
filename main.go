package main

import (
	"fmt"
	"os"
	"scriptprox/gui"
	"scriptprox/proxies/http"
	"scriptprox/proxies/udp"
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

	go updater.CheckForUpdates()

	if len(os.Args) > 1 {
		addrToConnectTo := os.Args[1]
		http.ChangeEndpoint(addrToConnectTo)
	}

	gui.RunGUI()
}
