# goMiddlewareChain [![Build Status](https://travis-ci.org/TobiEiss/goMiddlewareChain.svg?branch=master)](https://travis-ci.org/TobiEiss/goMiddlewareChain)

This is an express.js-like-middleware-chain for [julienschmidt's httprouter](https://github.com/julienschmidt/httprouter)

You can write your own middleware, and chain this to a lot of other middlewares (logging, auth,...).

## Getting started

### Install goMiddlewareChain
`go get github.com/TobiEiss/goMiddlewareChain`

### Your first API

Here a simple example with a simple Ping-Pong-Handler chained with a JSONResponseHandler (from templates).

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

## restricted-requestChainHandler
In some cases you need a restriction to apply requestChain. For example an auth-restriction.
You can use the `RestrictedRequestChainHandler`. If the `RestrictHandler` failed, the code doesn't pass the chain.

Same example with Auth:

```golang
package main

import (
	"log"
	"net/http"

	"github.com/TobiEiss/goMiddlewareChain"
	"github.com/TobiEiss/goMiddlewareChain/templates"
	"github.com/julienschmidt/httprouter"
)

// Ping return a simply pong
func Ping(response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) {
	// simply pong
	response.Status.Code = http.StatusOK
	response.Data = "pong"
}

func Auth(response *goMiddlewareChain.Response, request *http.Request, params httprouter.Params) bool {
	user := request.Header.Get("X-User")
	return user == "HomerSimpson"
}

func main() {
	router := httprouter.New()
	router.GET("/api/v0/ping", goMiddlewareChain.RestrictedRequestChainHandler(Auth, templates.JSONResponseHandler, Ping))

	log.Fatal(http.ListenAndServe(":8080", router))
}
```

Now run `curl --header "X-User: HomerSimpson" localhost:8080/api/v0/ping` in your terminal. You will get:
```json
{
    "data":"pong",
    "msg":"OK",
    "status":200
}
```

If you run `curl --header "X-User: BartSimpson" localhost:8080/api/v0/ping`, you get:
```json
{
	"msg":"failed by passing restrictHandler",
	"status":401
}
```

## handler from tamplates
- [logHandler](https://github.com/TobiEiss/goMiddlewareChain/blob/master/templates/logHandler.go) is an easy handler to log all accesses
- [jsonResponseHandler](https://github.com/TobiEiss/goMiddlewareChain/blob/master/templates/jsonResponseHandler.go) try to transform your response to valid json

You need more handler? Just let us now this.
