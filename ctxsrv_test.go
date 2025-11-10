package ctxsrv_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/koron-go/ctxsrv"
)

type dummyListener struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var _ net.Listener = (*dummyListener)(nil)

func (ln *dummyListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (ln *dummyListener) Close() error {
	ln.cancel()
	return nil
}

func (ln *dummyListener) Addr() net.Addr {
	return nil
}

func listenerFactory(ctx context.Context) func() (net.Listener, error) {
	return func() (net.Listener, error) {
		myCtx, cancel := context.WithCancel(ctx)
		return &dummyListener{ctx: myCtx, cancel: cancel}, nil
	}
}

func TestVerifyFailure(t *testing.T) {
	t.Run("Listen", func(t *testing.T) {
		cfg := &ctxsrv.Config{}
		err := cfg.ServeWithContext(t.Context())
		want := "no Listen function"
		if got := err.Error(); got != want {
			t.Errorf("unexpected error: got=%s", got)
		}
	})
	t.Run("Serve", func(t *testing.T) {
		cfg := &ctxsrv.Config{
			Listen: listenerFactory(t.Context()),
		}
		err := cfg.ServeWithContext(t.Context())
		want := "no Serve function"
		if got := err.Error(); got != want {
			t.Errorf("unexpected error: got=%s", got)
		}
	})
}

func TestServerImmediatelyTerminates(t *testing.T) {
	cfg := &ctxsrv.Config{
		Listen: listenerFactory(t.Context()),
		Serve:  func(net.Listener) error { return nil },
	}
	doneCalled := false
	cfg.WithDoneServer(func() {
		doneCalled = true
	})
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	err := cfg.ServeWithContext(ctx)
	if err != nil {
		t.Errorf("ServeWithContext returns an error: %s", err)
	}
	if !doneCalled {
		t.Error("DoneServer hook is not called")
	}
}
