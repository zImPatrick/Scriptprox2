package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"scriptprox/proxies/udp"
	"strings"
	"time"
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
		if addr != nil {
			stringAddr = strings.Split(addr.String(), ":")
		} else {
			stringAddr = []string{endpoint.Hostname(), fmt.Sprint(endpoint.Port())}
		}

		doScriptproxApiResponse(writer, 200, createScriptproxApiJSON(map[string]interface{}{
			"ip":   stringAddr[0],
			"port": stringAddr[1],
		}))
	case "/saveFiles":
		err := req.ParseMultipartForm(64 << 20)
		if err != nil {
			doScriptproxApiResponse(writer, 500, []byte("Couldn't parse form data"))
			return
		}
		receivedTime := fmt.Sprint(time.Now().UnixNano())
		folderPath := filepath.Join("exportedScripts", receivedTime)
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Couldn't create folder %s, %s\n", folderPath, err.Error())
			return
		}

		for _, fileHeaders := range req.MultipartForm.File {
			for _, header := range fileHeaders {
				fn := strings.ReplaceAll(header.Filename, "######", "/")
				savePath := filepath.Join(folderPath, fn)

				err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
				if err != nil {
					fmt.Printf("Couldn't create folder %s: %s\n", filepath.Dir(savePath), err.Error())
					continue
				}

				f, err := os.Create(savePath)
				if err != nil {
					fmt.Printf("Couldn't save %s, %s\n", savePath, err.Error())
					continue
				}
				defer f.Close()
				multipartFile, err := header.Open()
				if err != nil {
					fmt.Printf("Couldn't read contents of %s, %s\n", header.Filename, err.Error())
				}
				defer multipartFile.Close()
				io.Copy(f, multipartFile)
			}
		}
	default:
		{
			println(req.URL.Path)
			writer.WriteHeader(404)
		}
	}
}
