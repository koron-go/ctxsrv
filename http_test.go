package ctxsrv_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/koron-go/ctxsrv"
)

func TestHTTPContextCanceled(t *testing.T) {
	srv := &http.Server{
		Addr:    "127.0.0.1:",
		Handler: nil,
	}
	cfg := ctxsrv.HTTP(srv)
	doneCalled := false
	cfg.WithDoneContext(func() {
		doneCalled = true
	})
	ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
	defer cancel()
	err := cfg.ServeWithContext(ctx)
	if err != nil {
		t.Errorf("ServeWithContext failed: %s", err)
	}
	if !doneCalled {
		t.Error("DoneContext hook is not called")
	}
	http.ServeTLS()
}

func TestHTTPServerClose(t *testing.T) {
	srv := &http.Server{
		Addr:    "127.0.0.1:",
		Handler: nil,
	}
	cfg := ctxsrv.HTTP(srv)
	doneCalled := false
	cfg.WithDoneServer(func() {
		doneCalled = true
	})
	go func() {
		time.Sleep(200*time.Millisecond)
		srv.Close()
	}()
	err := cfg.ServeWithContext(t.Context())
	if err != nil {
		t.Errorf("ServeWithContext failed: %s", err)
	}
	if !doneCalled {
		t.Error("DoneServer hook is not called")
	}
}
