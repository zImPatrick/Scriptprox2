package main

import (
	"fmt"
	"os"
	"scriptprox/proxies/http"
	"scriptprox/proxies/udp"
	"sync"
)

func main() {
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)

	addrToConnectTo := os.Args[1]
	fmt.Println(addrToConnectTo)
	go http.InitProxy(&waitgroup)
	http.ChangeEndpoint(addrToConnectTo)
	go udp.InitProxy(&waitgroup)

	waitgroup.Wait()
}
