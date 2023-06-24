package interceptors

import (
	"net/http"
	"net/url"
	"scriptprox/settings"
)

func handlePreClientInterceptor(writer http.ResponseWriter, req *http.Request, reqBody *[]byte) *http.Request {
	parsedQuery, _ := url.ParseQuery(string(*reqBody))

	if parsedQuery.Get("method") == "initConnect" {
		if settings.GetSettings().GUIDOverride != "" && settings.GetSettings().TicketOverride != "" {
			parsedQuery.Set("guid", settings.GetSettings().GUIDOverride)
			parsedQuery.Set("cfxTicket", settings.GetSettings().TicketOverride)
		}
		encodedQuery := []byte(parsedQuery.Encode())
		*reqBody = encodedQuery
	}

	return req
}

var ClientPreInterceptor InterceptorBefore = InterceptorBefore{
	Endpoint: "/client",
	Method:   "POST",
	Handler:  handlePreClientInterceptor,
}
