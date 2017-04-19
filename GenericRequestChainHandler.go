package goMiddlewareChain

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GenericRequestChainHandler is the generic requestChainHandler.
//
// checkResponseOfEveryHandler 	-> if true: check the responseCode after every handler. If httpStatus is != 0 or 200, stop to move on.
// restrictHandler 				-> if != nil: pass first this handler. Use this for example for AuthHandler
// responseHandler 				-> if != nil: pass at last this handler. Use this for example for a JSON-converter
// handlers						-> all your handler-chain or middlewares or or or
func GenericRequestChainHandler(checkResponseOfEveryHandler bool, restrictHandler RestrictContextHandler, responseHandler ResponseHandler, handlers ...ContextHandler) httprouter.Handle {
	return httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		payload := Response{}
		rootContext := context.Background()
		var runningContext context.Context

		// check restriction
		allowed := true
		if restrictHandler != nil {
			runningContext, allowed = restrictHandler(rootContext, &payload, request, params)
		}

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

		// pass responseHandler
		if responseHandler != nil {
			responseHandler(&payload, writer, request, params)
		}
	})
}
