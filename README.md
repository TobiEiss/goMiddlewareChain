# goMiddlewareChain

This is a express.js like middleware chain for [julienschmidt's httprouter](github.com/julienschmidt/httprouter)

Stay tuned for more documentation

## Getting started

### Install goMiddlewareChain
`go get github.com/TobiEiss/goMiddlewareChain`

### Your first API

Here a simple example with a simple JSONResponseHandler.

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/TobiEiss/goMiddlewareChain"
	"github.com/TobiEiss/goMiddlewareChain/templates"
)

// Ping return a simply pong
func Ping(response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) {
	// simply pong
	response.Status.Code = http.StatusOK
	response.Data = "pong"
}

func main() {
	router := httprouter.New()
	router.GET("/api/v0/ping", goMiddlewareChain.RequestChainHandler(templates.JSONResponseHandler, Ping))

	log.Fatal(http.ListenAndServe(":8080", router))
}
```

After running this code, run `curl localhost:8080/api/v0/ping` in a terminal.
You will get the following:
```json
{
    "data":"pong",
    "msg":"OK",
    "status":200
}
```
Isn't it cool?

## TODOS
[ ] documentation   
[ ] tests   
[ ] travis