package goMiddlewareChain

import (
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
	bodyMap := map[string]interface{}{"data": response.Data}
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
