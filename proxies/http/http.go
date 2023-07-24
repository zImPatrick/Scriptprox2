package http

import (
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"scriptprox/interceptors"
	"scriptprox/utils"
	"strings"
	"sync"
)

var server *http.Server
var endpoint *url.URL
var httpClient *http.Client

var BeforeInterceptors []interceptors.InterceptorBefore = []interceptors.InterceptorBefore{
	interceptors.ClientPreInterceptor,
}
var AfterInterceptors []interceptors.InterceptorAfter = []interceptors.InterceptorAfter{
	interceptors.ClientInterceptor,
	interceptors.InfoInterceptor,
}

//go:embed server-tls.crt
var certPem []byte

//go:embed server-tls.key
var keyPem []byte

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
	if strings.HasPrefix(req.URL.Path, "/_api") {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/_api")
		handleScriptproxApiRequest(writer, req)
		return
	}
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

	for _, interceptor := range BeforeInterceptors {
		if interceptor.Endpoint == req.URL.Path && interceptor.Method == req.Method {
			interceptor.Handler(writer, req, &reqBody)
		}
	}

	fmt.Printf("%s %s %s\n", req.Method, req.URL.String(), reqBody)
	newReq, _ := http.NewRequest(
		req.Method,
		endpoint.JoinPath(
			req.URL.String(),
		).String(),
		strings.NewReader(string(reqBody)),
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

	// HTTP Client für Requests
	httpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// HTTP Server
	server = &http.Server{
		Addr:    ":30120",
		Handler: http.HandlerFunc(requestHandler),
	}
	go func () {
		err := server.ListenAndServe()
		if err != nil {
			utils.ThrowErrorAndQuit("Der HTTP-Server konnte nicht starten. Ist Port 30120 belegt?", err)
		}
	}()
	
	// HTTPS Server
	// du brauchst anscheinend einen HTTPS-Server
	// für die FiveM Resources. kp wieso
	keyPair, _ := tls.X509KeyPair(
		certPem,
		keyPem,
	)
	httpsServer := &http.Server{
		Addr:    ":30129",
		Handler: http.HandlerFunc(requestHandler),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{keyPair},
		},
	}
	err := httpsServer.ListenAndServeTLS("", "")
	if err != nil {
		utils.ThrowErrorAndQuit("Der HTTPS-Server konnte nicht starten. Ist Port 30129 belegt?", err)
	}
}

func ChangeEndpoint(endpointToChangeTo string) error {
	u, err := url.Parse(endpointToChangeTo)
	if err != nil {
		return errors.New("Ungültige Server-URL! " + err.Error())
	}
	endpoint = u
	return nil
}
