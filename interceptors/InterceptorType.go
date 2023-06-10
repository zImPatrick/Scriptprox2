package interceptors

import "net/http"

type InterceptorBefore struct {
	Endpoint string
	Method   string
	Handler  func(writer http.ResponseWriter, req *http.Request, reqBuffer []byte) *http.Request // returns updated request
}

type InterceptorAfter struct {
	Endpoint string
	Method   string
	Handler  func(writer http.ResponseWriter, req *http.Request, reqBuffer []byte, response *http.Response) *http.Response // returns updated response
}
