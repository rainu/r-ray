package http

import (
	"context"
	"github.com/rainu/r-ray/internal/http/controller"
	"net/http"
)

type server struct {
	server http.Server
}

func NewServer(addr string, headerPrefix string, processor controller.Processor) *server {
	return &server{
		server: http.Server{
			Addr:    addr,
			Handler: controller.NewProxy(headerPrefix, processor), //TODO: logging middleware
		},
	}
}

func (s *server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
