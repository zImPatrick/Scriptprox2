package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scriptprox/proxies/udp"
	"strings"
)

func doScriptproxApiResponse(writer http.ResponseWriter, statusCode int, body []byte) {
	writer.WriteHeader(statusCode)
	writer.Write(body)
}

func createScriptproxApiJSON(data map[string]interface{}) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return []byte(err.Error())
	}

	return jsonData
}

func handleScriptproxApiRequest(writer http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/mumble":
		addr := udp.GetAddr()
		var stringAddr []string
		if addr != nil { // das ist scheiße
			stringAddr = strings.Split(addr.String(), ":")
		} else {
			stringAddr = []string{endpoint.Hostname(), fmt.Sprint(endpoint.Port())}
		}

		doScriptproxApiResponse(writer, 200, createScriptproxApiJSON(map[string]interface{}{
			"ip":   stringAddr[0],
			"port": stringAddr[1],
		}))
	default:
		{
			println(req.URL.Path)
			writer.WriteHeader(404)
		}
	}
}
