package goMiddlewareChain

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"net/http/httptest"

	"encoding/json"

	"reflect"

	"github.com/julienschmidt/httprouter"
)

// response handler which write response to writer
var responseHandler = func(response *Response, writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	bodyMap := map[string]interface{}{"data": response.Data, "status": response.Status.Code}
	bodyJSON, _ := json.Marshal(bodyMap)
	fmt.Fprintf(writer, "%s", bodyJSON)
}

func TestMain(t *testing.T) {

	var tests = []struct {
		handler         Handler
		restrictHandler RestrictHandler
		expectedData    interface{}
	}{
		// Test 0:
		// Check if chain works
		{
			handler: func(response *Response, _ *http.Request, _ httprouter.Params) {
				response.Data = "data"
			},
			expectedData: "data",
		},
		// Test 1:
		// Check if restrictedChain works with restriction true
		{
			handler: func(response *Response, _ *http.Request, _ httprouter.Params) {
				response.Data = "data"
			},
			restrictHandler: func(response *Response, request *http.Request, _ httprouter.Params) bool {
				return true
			},
			expectedData: "data",
		},
	}

	// run all tests
	for index, test := range tests {
		// recorder and request to simulate test
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("POST", "", nil)

		// build chain
		if reflect.ValueOf(test.restrictHandler).IsNil() {
			handlerChain := RequestChainHandler(responseHandler, test.handler)
			handlerChain(recorder, request, nil)
		} else {
			restrictedHandleChain := RestrictedRequestChainHandler(test.restrictHandler, responseHandler, test.handler)
			restrictedHandleChain(recorder, request, nil)
		}

		// check result
		var responseBody map[string]interface{}
		json.Unmarshal([]byte(recorder.Body.String()), &responseBody)
		if err != nil || responseBody["data"] != test.expectedData {
			t.Errorf("Test %v failed: expected data is not equals response.data; expected: %v; response: %v", index, test.expectedData, responseBody["data"])
		}
	}
}

func TestContextHandler(t *testing.T) {

	var key = struct {
		key string
	}{
		key: "key",
	}

	var tests = []struct {
		handler         []ContextHandler
		restrictHandler RestrictHandler
		expectedData    interface{}
	}{
		// Test 0:
		// Check if chain works
		{
			handler: []ContextHandler{
				func(ctx context.Context, response *Response, _ *http.Request, _ httprouter.Params) context.Context {
					return context.WithValue(ctx, key, "data")
				},
				func(ctx context.Context, response *Response, _ *http.Request, _ httprouter.Params) context.Context {
					response.Data = ctx.Value(key).(string)
					return ctx
				},
			},
			expectedData: "data",
		},
	}

	// run all tests
	for index, test := range tests {
		// recorder and request to simulate test
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("POST", "", nil)

		// build chain
		handlerChain := RequestChainContextHandler(responseHandler, test.handler...)
		handlerChain(recorder, request, nil)

		// check result
		var responseBody map[string]interface{}
		json.Unmarshal([]byte(recorder.Body.String()), &responseBody)
		if err != nil || responseBody["data"] != test.expectedData {
			t.Errorf("Test %v failed: expected data is not equals response.data; expected: %v; response: %v", index, test.expectedData, responseBody["data"])
		}
	}
}

func TestContextFailHandler(t *testing.T) {
	var tests = []struct {
		handler                     []ContextHandler
		restrictHandler             RestrictContextHandler
		expectedStatus              int
		checkResponseOfEveryHandler bool
	}{
		// Test 0:
		// Check expected code is 500
		{
			handler: []ContextHandler{
				func(ctx context.Context, response *Response, _ *http.Request, _ httprouter.Params) context.Context {
					response.Status.Code = http.StatusOK
					return ctx
				},
				func(ctx context.Context, response *Response, _ *http.Request, _ httprouter.Params) context.Context {
					response.Status.Code = http.StatusInternalServerError
					return ctx
				},
				// This handler should not run!
				func(ctx context.Context, response *Response, _ *http.Request, _ httprouter.Params) context.Context {
					response.Status.Code = http.StatusOK
					return ctx
				},
			},
			restrictHandler: func(ctx context.Context, _ *Response, _ *http.Request, _ httprouter.Params) (context.Context, bool) {
				return ctx, true
			},
			expectedStatus:              500,
			checkResponseOfEveryHandler: true,
		},

		// Test 1:
		// Check expected code == 200
		{
			handler: []ContextHandler{
				func(ctx context.Context, response *Response, _ *http.Request, _ httprouter.Params) context.Context {
					response.Status.Code = http.StatusOK
					return ctx
				},
				func(ctx context.Context, response *Response, _ *http.Request, _ httprouter.Params) context.Context {
					response.Status.Code = http.StatusInternalServerError
					return ctx
				},
				// This handler should not run!
				func(ctx context.Context, response *Response, _ *http.Request, _ httprouter.Params) context.Context {
					response.Status.Code = http.StatusOK
					return ctx
				},
			},
			restrictHandler: func(ctx context.Context, _ *Response, _ *http.Request, _ httprouter.Params) (context.Context, bool) {
				return ctx, true
			},
			expectedStatus:              200,
			checkResponseOfEveryHandler: false,
		},
	}

	// run all tests
	for index, test := range tests {
		// recorder and request to simulate test
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("POST", "", nil)

		// build chain
		handlerChain := RestrictedRequestChainContextHandlerWithResponseCheck(test.checkResponseOfEveryHandler, test.restrictHandler, responseHandler, test.handler...)
		handlerChain(recorder, request, nil)

		// check result
		var responseBody map[string]interface{}
		json.Unmarshal([]byte(recorder.Body.String()), &responseBody)
		if err != nil || int(responseBody["status"].(float64)) != test.expectedStatus {
			t.Errorf("Test %v failed: expected data is not equals response.status.code; expected: %v; response: %v", index, test.expectedStatus, responseBody["status"])
		}
	}
}
