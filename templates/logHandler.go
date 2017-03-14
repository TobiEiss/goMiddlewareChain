package templates

import (
	"context"
	"net/http"

	"log"

	"github.com/TobiEiss/goMiddlewareChain"
	"github.com/julienschmidt/httprouter"
)

// LogHandler log requests
func LogHandler(response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) {
	log.Println(request.RemoteAddr, "-", request.Method, "-", request.Host, "-", request.URL, "-", request.Header)
}

// LogContextHandler log requests
func LogContextHandler(ctx context.Context, response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) context.Context {
	LogHandler(response, request, params)
	return ctx
}
