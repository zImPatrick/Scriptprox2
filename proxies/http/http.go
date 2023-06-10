package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"scriptprox/interceptors"
	"sync"
)

var server *http.Server
var endpoint *url.URL
var httpClient *http.Client

var BeforeInterceptors []interceptors.InterceptorBefore = []interceptors.InterceptorBefore{}
var AfterInterceptors []interceptors.InterceptorAfter = []interceptors.InterceptorAfter{
	interceptors.ClientInterceptor,
}

func TransferHeaders(to http.ResponseWriter, from http.Header) {
	for name, header := range from {
		for _, head := range header {
			to.Header().Add(name, head)
		}
	}
}

func WriteError(writer http.ResponseWriter, err string) {
	writer.WriteHeader(500)
	encoded, _ := json.Marshal(map[string]interface{}{
		"error": fmt.Sprintf("Scriptprox: %s", err),
	})
	writer.Write([]byte(encoded))
}

func requestHandler(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/spawnmanager/resource.rpf" {
		f, err := os.Open("resource.rpf")
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("Konnte resource.rpf nicht öffnen"))
		}
		defer f.Close()

		writer.WriteHeader(200)
		io.Copy(writer, f)
		return
	}

	if endpoint == (&url.URL{}) {
		writer.WriteHeader(500)
		writer.Write([]byte("no endpoint given"))

		return
	}

	// die req body ist meist ziemlich klein, und aus convienience gründen lesen wir einfach den body sofort
	reqBody, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

	newReq, _ := http.NewRequest(
		req.Method,
		endpoint.JoinPath(
			req.URL.String(),
		).String(),
		req.Body,
	)
	newReq.Header = req.Header
	newResp, err := httpClient.Do(newReq)
	if err != nil {
		WriteError(writer, "Fehler beim Verbinden: "+err.Error())
		return
	}

	for _, interceptor := range AfterInterceptors {
		if interceptor.Endpoint == req.URL.Path && interceptor.Method == req.Method {
			interceptor.Handler(writer, req, reqBody, newResp)
		}
	}

	TransferHeaders(writer, newResp.Header)
	writer.WriteHeader(newResp.StatusCode)
	flushCopy(writer, newResp.Body)
}

func InitProxy(waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	httpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	server = &http.Server{
		Addr:    ":30120",
		Handler: http.HandlerFunc(requestHandler),
	}
	httpsServer := &http.Server{
		Addr:    ":30129",
		Handler: http.HandlerFunc(requestHandler),
	}

	go server.ListenAndServe()
	httpsServer.ListenAndServeTLS("server-tls.crt", "server-tls.key")
}

func ChangeEndpoint(endpointToChangeTo string) {
	u, err := url.Parse(endpointToChangeTo)
	if err != nil {
		panic(err)
	}
	endpoint = u
}
