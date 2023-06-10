package interceptors

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"scriptprox/proxies/udp"

	"golang.org/x/exp/slices"
)

func handleClientInterceptor(writer http.ResponseWriter, req *http.Request, reqBody []byte, resp *http.Response) *http.Response {
	// configure udp client
	parsedQuery, _ := url.ParseQuery(string(reqBody))

	if parsedQuery.Get("method") == "getEndpoints" {
		var addrToConnectTo *net.UDPAddr
		readBody, _ := io.ReadAll(resp.Body)
		var bodyJson []string
		json.Unmarshal(readBody, &bodyJson)

		if len(bodyJson) > 0 { // ein endpoint ist angegeben
			addrToConnectTo, _ = net.ResolveUDPAddr("udp", bodyJson[0])
		} else { // default endpoint nehmen
			addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(resp.Request.URL.Hostname(), resp.Request.URL.Port()))

			if err != nil {
				fmt.Println(net.JoinHostPort(resp.Request.URL.Hostname(), resp.Request.URL.Port()))
				fmt.Println(err)
			}
			addrToConnectTo = addr
		}

		udp.ConnectToUDP(addrToConnectTo)

		fmt.Println(bodyJson)

		resp.Body = io.NopCloser(bytes.NewBuffer(
			[]byte("[\"127.0.0.1:30120\"]"),
		))
	}

	if parsedQuery.Get("method") == "getConfiguration" {
		readBody, _ := io.ReadAll(resp.Body)
		var bodyData GetConfigurationData
		err := json.Unmarshal(readBody, &bodyData)
		if err != nil {
			fmt.Println("err" + err.Error())
		}

		idx := slices.IndexFunc(bodyData.Resources, func(res ResourceData) bool {
			return res.Name == "spawnmanager"
		})

		if idx != -1 {
			bodyData.Resources = slices.Delete(bodyData.Resources, idx, idx+1)
		} else {
			idx = 0
		}

		for ind, resource := range bodyData.Resources {
			if resource.FileServer == "" {
				fmt.Println(resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + "/files")
				bodyData.Resources[ind].FileServer = resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + "/files"
			}
		}

		bodyData.Resources = slices.Insert(bodyData.Resources, idx, ResourceData{
			Name: "spawnmanager",
			Files: map[string]string{
				"resource.rpf": hashResourceRPF(),
			},
			StreamFiles: map[string]string{},
			FileServer:  "https://127.0.0.1:30129",
		})

		marshaled, _ := json.Marshal(bodyData)
		resp.Body = io.NopCloser(bytes.NewBuffer(marshaled))

		// resp.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}
	return resp
}

func hashResourceRPF() string {
	f, err := os.Open("resource.rpf")
	if err != nil {
		fmt.Println("Ein Fehler ist beim Öffnen von resource.rpf aufgetreten. Hast du resource.rpf im Ordner?")
	}
	defer f.Close()

	h := sha1.New()
	io.Copy(h, f)

	return fmt.Sprintf("%x", h.Sum(nil))
}

var ClientInterceptor InterceptorAfter = InterceptorAfter{
	Endpoint: "/client",
	Method:   "POST",
	Handler:  handleClientInterceptor,
}

type ResourceData struct {
	Name        string            `json:"name"`
	Files       map[string]string `json:"files"`
	StreamFiles any               `json:"streamFiles"`
	FileServer  string            `json:"fileServer,omitempty"`
	// StreamFiles StreamFiles    `json:"streamFiles"`
	// streamFiles haben andere Properties und sind uns ziemlich egal
	URI string `json:"uri,omitempty"`
}

type GetConfigurationData struct {
	FileServer string         `json:"fileServer"`
	Resources  []ResourceData `json:"resources"`
}
