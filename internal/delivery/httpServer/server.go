package httpServer

import (
	"eljur/internal/config"
	"net"
	"net/http"
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
