package goMiddlewareChain

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RequestChainHandler chains all handler
func RequestChainHandler(responseHandler ResponseHandler, handlers ...Handler) httprouter.Handle {
	return httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		payload := Response{}

		// iterate all handlers
		for _, handler := range handlers {
			handler(&payload, request, params)
		}

		// pass responseHandler
		responseHandler(&payload, writer, request, params)
	})
}

// RestrictedRequestChainHandler need a RestrictHandler.
// A RestrictHandler returns bool if call is allowed.
func RestrictedRequestChainHandler(restrictHandler RestrictHandler, responseHandler ResponseHandler, handlers ...Handler) httprouter.Handle {
	return httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		payload := Response{}

		// check restriction
		allowed := restrictHandler(&payload, request, params)

		if allowed {
			// iterate all handlers
			for _, handler := range handlers {
				handler(&payload, request, params)
			}
		}

		// pass ResponseHandler
		responseHandler(&payload, writer, request, params)
	})
}
