package interceptors

import "net/http"

type InterceptorBefore struct {
	endpoint string
	method   string
	Handler  func(writer http.ResponseWriter, req *http.Request) *http.Request // returns updated request
}

type InterceptorAfter struct {
	endpoint string
	method   string
	Handler  func(writer http.ResponseWriter, req *http.Request, response *http.Response) *http.Response // returns updated response
}
