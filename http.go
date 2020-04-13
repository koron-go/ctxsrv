package ctxsrv

import (
	"net"
	"net/http"
)

// HTTP creates a config for HTTP serer.
func HTTP(srv *http.Server) *Config {
	return &Config{
		Listen: func() (net.Listener, error) {
			addr := srv.Addr
			if addr == "" {
				addr = ":http"
			}
			return net.Listen("tcp", addr)
		},
		Serve:    srv.Serve,
		Shutdown: srv.Shutdown,
	}
}

// HTTPS creates a config for HTTPS(TLS) server.
func HTTPS(srv *http.Server, certFile, keyFile string) *Config {
	return &Config{
		Listen: func() (net.Listener, error) {
			addr := srv.Addr
			if addr == "" {
				addr = ":https"
			}
			return net.Listen("tcp", addr)
		},
		Serve: func(ln net.Listener) error {
			return srv.ServeTLS(ln, certFile, keyFile)
		},
		Shutdown: srv.Shutdown,
	}
}
