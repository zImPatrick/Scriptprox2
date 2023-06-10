package http

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
)

var server *http.Server
var endpoint *url.URL
var httpClient *http.Client

func TransferHeaders(to http.ResponseWriter, from http.Header) {
	for name, header := range from {
		for _, head := range header {
			to.Header().Add(name, head)
		}
	}
}

func requestHandler(writer http.ResponseWriter, req *http.Request) {
	if endpoint == (&url.URL{}) {
		writer.WriteHeader(500)
		writer.Write([]byte("no endpoint given"))

		return
	}

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
		println("error" + err.Error())
		return
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

	server.ListenAndServe()
}

func ChangeEndpoint(endpointToChangeTo string) {
	u, err := url.Parse(endpointToChangeTo)
	if err != nil {
		panic(err)
	}
	endpoint = u
}
