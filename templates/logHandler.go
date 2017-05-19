package templates

import (
	"context"
	"log/syslog"
	"net/http"

	"log"

	"github.com/TobiEiss/goMiddlewareChain"
	"github.com/julienschmidt/httprouter"
)

// LogHandler log requests.
// Use this handler to log every call to console.
func LogHandler(response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) {
	log.Println(request.RemoteAddr, "-", request.Method, "-", request.Host, "-", request.URL, "-", request.Header)
}

// LogContextHandler log requests
// Use this context-handler to log every call to console.
func LogContextHandler(ctx context.Context, response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) context.Context {
	LogHandler(response, request, params)
	return ctx
}

// LoggerHandler holds the logger-instance
type LoggerHandler struct {
	Writer *syslog.Writer
}

// NewLoggerHandler creates a new LoggerHandler
func NewLoggerHandler(writer *syslog.Writer) *LoggerHandler {
	return &LoggerHandler{Writer: writer}
}

// LoggerContextKey a key to map the logger
var LoggerContextKey = goMiddlewareChain.ContextKey{Key: ""}

// Handle is the handler for the LoggerHandler
func (handler *LoggerHandler) Handle(ctx context.Context, response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) context.Context {
	return context.WithValue(ctx, LoggerContextKey, handler)
}
