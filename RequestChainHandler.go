package goMiddlewareChain

import (
	"net/http"

	"digital.fino/contractsafe-api/middleware"
	"github.com/julienschmidt/httprouter"
)

// RequestChainHandler chains all handler
func RequestChainHandler(responseHandler middleware.ResponseHandler, handlers ...middleware.Handler) httprouter.Handle {
	return httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		payload := middleware.Response{}

		// iterate all handlers
		for _, handler := range handlers {
			handler(&payload, request, params)
		}

		// pass responseHandler
		responseHandler(&payload, writer, request, params)
	})
}
