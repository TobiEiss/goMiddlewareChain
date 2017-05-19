package goMiddlewareChain

import (
	"context"
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

// ContextHandler a handler with the go-lang-context
type ContextHandler func(context.Context, *Response, *http.Request, httprouter.Params) context.Context

// RestrictContextHandler restrict handler for contextHandler
type RestrictContextHandler func(context.Context, *Response, *http.Request, httprouter.Params) (context.Context, bool)

// ContextKey to map ContextValues
type ContextKey struct {
	Key string
}
