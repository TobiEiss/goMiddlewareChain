package goMiddlewareChain

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RequestChainHandler chains all handler
func RequestChainHandler(responseHandler ResponseHandler, handlers ...Handler) httprouter.Handle {
	return RequestChainHandlerWithResponseCheck(false, responseHandler, handlers...)
}

// RequestChainHandlerWithResponseCheck chains all handler and check every response
func RequestChainHandlerWithResponseCheck(checkResponseOfEveryHandler bool, responseHandler ResponseHandler, handlers ...Handler) httprouter.Handle {
	return httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		payload := Response{}

		// iterate all handlers
		for _, handler := range handlers {
			handler(&payload, request, params)
			if checkResponseOfEveryHandler && (payload.Status.Code != http.StatusOK && payload.Status.Code != 0) {
				break
			}
		}

		// pass responseHandler
		responseHandler(&payload, writer, request, params)
	})
}

// RequestChainContextHandler chains all handler
func RequestChainContextHandler(responseHandler ResponseHandler, handlers ...ContextHandler) httprouter.Handle {
	return httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		payload := Response{}
		rootContext := context.Background()

		// iterate all handlers
		var runningContext context.Context
		runningContext = rootContext
		for _, handler := range handlers {
			runningContext = handler(runningContext, &payload, request, params)
		}

		// pass responseHandler
		responseHandler(&payload, writer, request, params)
	})
}

// RestrictedRequestChainHandler need a RestrictHandler.
// A RestrictHandler returns bool if call is allowed.
func RestrictedRequestChainHandler(restrictHandler RestrictHandler, responseHandler ResponseHandler, handlers ...Handler) httprouter.Handle {
	return RestrictedRequestChainHandlerWithResponseCheck(false, restrictHandler, responseHandler, handlers...)
}

// RestrictedRequestChainHandlerWithResponseCheck need a RestrictHandler.
// If checkResponseOfEveryHandler is true, handler check every response.
func RestrictedRequestChainHandlerWithResponseCheck(checkResponseOfEveryHandler bool, restrictHandler RestrictHandler, responseHandler ResponseHandler, handlers ...Handler) httprouter.Handle {
	return httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		payload := Response{}

		// check restriction
		allowed := restrictHandler(&payload, request, params)

		if allowed {
			// iterate all handlers
			for _, handler := range handlers {
				handler(&payload, request, params)
				if checkResponseOfEveryHandler && (payload.Status.Code != http.StatusOK && payload.Status.Code != 0) {
					break
				}
			}
		} else if payload.Status.Code == 0 {
			payload.Status.Code = http.StatusUnauthorized
			payload.Status.Message = "failed by passing restrictHandler"
		}

		// pass ResponseHandler
		responseHandler(&payload, writer, request, params)
	})
}

// RestrictedRequestChainContextHandler need a RestrictHandler.
// A RestrictHandler returns bool if call is allowed.
func RestrictedRequestChainContextHandler(restrictHandler RestrictContextHandler, responseHandler ResponseHandler, handlers ...ContextHandler) httprouter.Handle {
	return RestrictedRequestChainContextHandlerWithResponseCheck(false, restrictHandler, responseHandler, handlers...)
}

// RestrictedRequestChainContextHandlerWithResponseCheck exec all handlers
// If checkResponseOfEveryHandler is true, handler check every response.
func RestrictedRequestChainContextHandlerWithResponseCheck(checkResponseOfEveryHandler bool, restrictHandler RestrictContextHandler, responseHandler ResponseHandler, handlers ...ContextHandler) httprouter.Handle {
	return httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		payload := Response{}
		rootContext := context.Background()
		var runningContext context.Context

		// check restriction
		runningContext, allowed := restrictHandler(rootContext, &payload, request, params)

		if allowed {
			// iterate all handlers
			runningContext = rootContext
			for _, handler := range handlers {
				runningContext = handler(runningContext, &payload, request, params)
				if checkResponseOfEveryHandler && (payload.Status.Code != http.StatusOK && payload.Status.Code != 0) {
					break
				}
			}
		} else if payload.Status.Code == 0 {
			payload.Status.Code = http.StatusUnauthorized
			payload.Status.Message = "failed by passing restrictHandler"
		}

		// pass ResponseHandler
		responseHandler(&payload, writer, request, params)
	})
}
