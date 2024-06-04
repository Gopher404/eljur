package httpServer

import (
	"context"
	"eljur/internal/config"
	"net"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func NewServer(handler http.Handler, cnf *config.BindConfig) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:    net.JoinHostPort(cnf.Ip, cnf.Port),
			Handler: handler,
		},
	}
}

func (s *HttpServer) Run() error {
	return s.server.ListenAndServe()
}

func (s *HttpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return s.server.Shutdown(ctx)
}
