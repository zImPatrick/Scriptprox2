package main

import (
	"os"
	"scriptprox/proxies/http"
	"scriptprox/proxies/udp"
	"sync"
)

func main() {
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)
	go http.InitProxy(&waitgroup)
	go udp.InitProxy(&waitgroup)

	if len(os.Args) > 0 {
		addrToConnectTo := os.Args[1]
		http.ChangeEndpoint(addrToConnectTo)
	}

	go cli()
	waitgroup.Wait()
}
