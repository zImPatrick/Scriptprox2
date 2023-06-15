package interceptors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scriptprox/settings"
)

func handleInfoInterceptor(writer http.ResponseWriter, req *http.Request, reqBody []byte, resp *http.Response) *http.Response {
	fmt.Println("info.json")
	infoData, _ := io.ReadAll(resp.Body)
	var contents InfoEndpointContent
	json.Unmarshal(infoData, &contents) // todo: handle error
	if settings.GetSettings().HWIDBypass {
		delete(contents.Vars, "sv_licenseKeyToken")
	}

	marshaled, _ := json.Marshal(contents)
	resp.Header.Set("Content-Length", fmt.Sprint(len(marshaled)))
	resp.Body = io.NopCloser(bytes.NewBuffer(marshaled))
	return resp
}

var InfoInterceptor InterceptorAfter = InterceptorAfter{
	Endpoint: "/info.json",
	Method:   "GET",
	Handler:  handleInfoInterceptor,
}

type InfoEndpointContent struct {
	EnhancedHostSupport string            `json:"enhancedHostSupport"`
	RequestSteamTicket  string            `json:"requestSteamTicket"`
	Resources           []string          `json:"resources"`
	Server              string            `json:"server"`
	Vars                map[string]string `json:"vars"`
	Version             int               `json:"version"`
}
