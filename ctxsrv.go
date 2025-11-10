// Package ctxsrv provides context aware `Serve()` for net/http
package ctxsrv

import (
	"context"
	"errors"
	"net"
	"time"
)

// ListenFunc provides constructor of net.Listener for a service.
type ListenFunc func() (net.Listener, error)

// ServeFunc provides a service (ex. HTTP).
type ServeFunc func(net.Listener) error

// ShutdownFunc provides a function to shutdown a service.
type ShutdownFunc func(context.Context) error

// DoneFunc provides function to monitor `Done` event.
type DoneFunc func()

// Config provides configuration for `ctxsrv.Serve()` function.
type Config struct {
	Listen   ListenFunc
	Serve    ServeFunc
	Shutdown ShutdownFunc

	DoneContext DoneFunc
	DoneServer  DoneFunc

	// ShutdownTimeout specifies duration for timeout of shutdown.
	ShutdownTimeout time.Duration
}

// WithDoneContext sets the DoneContext handler.
func (cfg *Config) WithDoneContext(fn DoneFunc) *Config {
	cfg.DoneContext = fn
	return cfg
}

// WithDoneServer sets the DoneServer handler.
func (cfg *Config) WithDoneServer(fn DoneFunc) *Config {
	cfg.DoneServer = fn
	return cfg
}

// WithShutdownTimeout sets the timeout duration to shutdown.
func (cfg *Config) WithShutdownTimeout(d time.Duration) *Config {
	cfg.ShutdownTimeout = d
	return cfg
}

// ServeWithContext just calls with `ctxsrv.Serve(ctx, *cfg)`
func (cfg *Config) ServeWithContext(ctx context.Context) error {
	return Serve(ctx, *cfg)
}

func (cfg Config) verify() error {
	if cfg.Listen == nil {
		return errors.New("no Listen function")
	}
	if cfg.Serve == nil {
		return errors.New("no Serve function")
	}
	return nil
}

func (cfg Config) shutdown() error {
	if cfg.Shutdown == nil {
		return nil
	}
	ctx := context.Background()
	if cfg.ShutdownTimeout > 0 {
		ctx2, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
		defer cancel()
		ctx = ctx2
	}
	return cfg.Shutdown(ctx)
}

// Serve serves a service based on context. It is shutdown if the context is
// done.
func Serve(ctx context.Context, cfg Config) error {
	if err := cfg.verify(); err != nil {
		return err
	}
	ln, err := cfg.Listen()
	if err != nil {
		return err
	}
	defer ln.Close()
	srvCtx, srvCancel := context.WithCancel(context.Background())
	defer srvCancel()
	ch := make(chan error)
	go func() {
		select {
		case <-ctx.Done():
			if cfg.DoneContext != nil {
				cfg.DoneContext()
			}
			ch <- cfg.shutdown()
		case <-srvCtx.Done():
			if cfg.DoneServer != nil {
				cfg.DoneServer()
			}
			ch <- nil
		}
		close(ch)
	}()
	err = cfg.Serve(ln)
	srvCancel()
	if err != nil {
		<-ch
		return err
	}
	return <-ch
}
