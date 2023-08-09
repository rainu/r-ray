package http

import (
	"context"
	"net/http"
)

type server struct {
	server http.Server
}

func NewServer(addr string, handler http.Handler) *server {
	return &server{
		server: http.Server{
			Addr:    addr,
			Handler: handler, //TODO: logging middleware
		},
	}
}

func (s *server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
