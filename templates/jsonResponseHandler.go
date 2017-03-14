package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TobiEiss/goMiddlewareChain"
	"github.com/julienschmidt/httprouter"
)

// JSONContentType The default content-type
const JSONContentType = "application/json"

// NotFoundResponseHandler handle the default 404 errors with JSONResponseHandler
func NotFoundResponseHandler(writer http.ResponseWriter, request *http.Request) {
	JSONResponseHandler(&goMiddlewareChain.Response{Status: goMiddlewareChain.Status{Code: http.StatusNotFound}}, writer, request, nil)
}

// MethodNotAllowedResponseHandler handle the default 405 errors with JSONResponseHandler
func MethodNotAllowedResponseHandler(writer http.ResponseWriter, request *http.Request) {
	JSONResponseHandler(&goMiddlewareChain.Response{Status: goMiddlewareChain.Status{Code: http.StatusMethodNotAllowed}}, writer, request, nil)
}

// PanicHandler handle all crashes with a proper JSONResponseHandler response
func PanicHandler(writer http.ResponseWriter, request *http.Request, p interface{}) {
	JSONResponseHandler(&goMiddlewareChain.Response{Status: goMiddlewareChain.Status{Code: http.StatusInternalServerError}}, writer, request, nil)
}

// JSONResponseHandler wrap the response in a standard json structure
func JSONResponseHandler(response *goMiddlewareChain.Response, writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	if response.Status.Code == 0 {
		response.Status.Code = http.StatusInternalServerError
	}

	if response.Status.Message == "" {
		response.Status.Message = http.StatusText(response.Status.Code)
	}

	// set proper response headers
	writer.Header().Set("Content-Type", JSONContentType)
	writer.WriteHeader(int(response.Status.Code))

	// create default struct
	body := make(map[string]interface{})

	// set response body
	body["status"] = response.Status.Code
	body["msg"] = response.Status.Message

	if response.Data != nil {
		body["data"] = response.Data
	}

	// marshal to json and send
	responseByte, err := json.Marshal(body)
	if err != nil {
		log.Println("Can't marshal response-body", err)
	}

	fmt.Fprintf(writer, "%s", responseByte)
}

// JSONResponseContextHandler wrap the response in a standard json structure
func JSONResponseContextHandler(ctx context.Context, response *goMiddlewareChain.Response, writer http.ResponseWriter, request *http.Request, params httprouter.Params) context.Context {
	JSONResponseHandler(response, writer, request, params)
	return ctx
}
