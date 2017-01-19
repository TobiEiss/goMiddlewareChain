package goMiddlewareChain

import (
	"fmt"
	"net/http"
	"testing"

	"net/http/httptest"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
)

func TestMain(t *testing.T) {

	var tests = []struct {
		handler      Handler
		expectedData interface{}
	}{
		// Test 0:
		{
			handler: func(response *Response, _ *http.Request, _ httprouter.Params) {
				response.Data = "data"
			},
			expectedData: "data",
		},
	}

	// response handler which write response to writer
	responseHandler := func(response *Response, writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		bodyMap := map[string]interface{}{"data": response.Data}
		bodyJSON, _ := json.Marshal(bodyMap)
		fmt.Fprintf(writer, "%s", bodyJSON)
	}

	// run all tests
	for index, test := range tests {
		// recorder and request to simulate test
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "", nil)

		// build chain
		handlerChain := RequestChainHandler(responseHandler, test.handler)
		handlerChain(recorder, request, nil)

		// check result
		var responseBody map[string]interface{}
		json.Unmarshal([]byte(recorder.Body.String()), &responseBody)
		if err != nil || responseBody["data"] != test.expectedData {
			t.Errorf("Test %v failed: expected data is not equals response.data; expected: %v; response: %v", index, test.expectedData, responseBody["data"])
		}
	}
}
