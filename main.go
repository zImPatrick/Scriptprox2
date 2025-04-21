package main

import (
	"scriptprox/gui"
	"scriptprox/proxies/http"
	"scriptprox/proxies/udp"
	"sync"
)

func main() {
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)
	go http.InitProxy(&waitgroup)
	go udp.InitProxy(&waitgroup)

	gui.RunGUI()
}
