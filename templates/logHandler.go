package templates

import (
	"net/http"

	"log"

	"github.com/TobiEiss/goMiddlewareChain"
	"github.com/julienschmidt/httprouter"
)

// LogHandler log requests
func LogHandler(response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) {
	log.Println(request.RemoteAddr, "-", request.Method, "-", request.Host, "-", request.URL, "-", request.Header)
}
