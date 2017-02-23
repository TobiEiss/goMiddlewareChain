package templates

import (
	"net/http"

	"github.com/TobiEiss/goMiddlewareChain"
	"github.com/TobiEiss/httprouter"
)

// Ping route should just respond with pong
func Ping(response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) {
	response.Data = "PONG"
	response.Status.Code = http.StatusOK
}
