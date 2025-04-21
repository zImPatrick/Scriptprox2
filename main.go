package main

import (
	"scriptprox/gui"
	"scriptprox/proxies/http"
	"scriptprox/proxies/udp"
)

func main() {
	go http.InitProxy()
	go udp.InitProxy()
	gui.RunGUI()
}
