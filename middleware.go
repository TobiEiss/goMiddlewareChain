package goMiddlewareChain

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Status represents the handler status
type Status struct {
	Code    int
	Message string
}

// Response struct for Handlers
type Response struct {
	Data   interface{}
	Status Status
}

// Handler represent a chainable Handler (middleware-like)
type Handler func(*Response, *http.Request, httprouter.Params)

// RestrictHandler restricts to handle following handlers
type RestrictHandler func(*Response, *http.Request, httprouter.Params) bool

// ResponseHandler required for every Endpoint
type ResponseHandler func(*Response, http.ResponseWriter, *http.Request, httprouter.Params)
