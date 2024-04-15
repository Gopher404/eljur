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

var timeOut = time.Second * 10

func NewServer(handler http.Handler, cnf *config.BindConfig) *HttpServer {
	timeOut = cnf.TimeOut
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
