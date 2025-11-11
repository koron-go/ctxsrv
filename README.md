# koron-go/ctxsrv

[![GoDoc](https://godoc.org/github.com/koron-go/ctxsrv?status.svg)](https://godoc.org/github.com/koron-go/ctxsrv)
[![Actions/Go](https://github.com/koron-go/ctxsrv/workflows/Go/badge.svg)](https://github.com/koron-go/ctxsrv/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/ctxsrv)](https://goreportcard.com/report/github.com/koron-go/ctxsrv)

## Overview

`ctxsrv` is a Go package which provides a `context.Context` aware server.
It helps to implement graceful shutdown mechanism for your servers.

## Getting Started

Install or update the package:

```console
$ go get github.com/koron-go/ctxsrv@latest
```

## Usage

To use `ctxsrv`, you need to create a `ctxsrv.Config` object and pass it to the `ctxsrv.Serve` function. The `Config` object requires three functions:

*   `ListenFunc`: Creates a `net.Listener`.
*   `ServeFunc`: Starts the server with the created `net.Listener`.
*   `ShutdownFunc`: (OPTIONAL) Shuts down the server gracefully.

The `ctxsrv` package provides helper functions to create a `Config` for `net/http` servers:

*   `ctxsrv.HTTP`: Creates a `Config` for an HTTP server.
*   `ctxsrv.HTTPS`: Creates a `Config` for an HTTPS server.

## Sample Code

Here is an example of an HTTP server with graceful shutdown:

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/koron-go/ctxsrv"
)

func main() {
	// Create a context that is canceled on an interrupt signal.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Create an HTTP server.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server with graceful shutdown.
	err := ctxsrv.HTTP(srv).
		WithShutdownTimeout(5 * time.Second).
		ServeWithContext(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```
