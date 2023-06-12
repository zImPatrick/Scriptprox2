package main

import (
	"fmt"
	"os"
	"scriptprox/gui"
	"scriptprox/proxies/http"
	"scriptprox/proxies/udp"
	"sync"
)

func main() {
	fmt.Println("Starting....")
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)
	go http.InitProxy(&waitgroup)
	go udp.InitProxy(&waitgroup)

	if len(os.Args) > 1 {
		addrToConnectTo := os.Args[1]
		http.ChangeEndpoint(addrToConnectTo)
	}

	go cli()
	go gui.RunGUI()
	waitgroup.Wait()
}
