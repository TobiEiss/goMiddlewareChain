# goMiddlewareChain

This is a express.js-like-middleware-chain for [julienschmidt's httprouter](https://github.com/julienschmidt/httprouter)

You can write your own middleware, and chain this to a lot of other middlewares (logging, auth,...).

## Getting started

### Install goMiddlewareChain
`go get github.com/TobiEiss/goMiddlewareChain`

### Your first API

Here a simple example with a simple Ping-Pong-Hanlder chained with a JSONResponseHandler (from templates).

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

## handler from tamplates
- [logHandler](loghandler) is an easy handler to log all accesses
- jsonResponseHandler try to transform your response to valid json

You need more handler? Just let me now this and write me.

## TODOS
[ ] tests   
[ ] travis