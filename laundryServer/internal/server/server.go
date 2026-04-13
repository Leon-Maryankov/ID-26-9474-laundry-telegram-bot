package server

import (
	"net/http"

	"laundryServer/internal/multiplexer"
	"laundryServer/internal/slots"
	"laundryServer/internal/users"
)

type Server struct {
	server *http.Server
}

func NewServer(host, token string, users *users.Users, slots *slots.Slots) (*Server, error) {
	result := &Server{
		server: &http.Server{
			Addr: host,
			Handler: multiplexer.NewMultiplexer(token, users, slots),
		},
	}

	return result, nil
}

func (server *Server) ListenAndServe() {
	server.server.ListenAndServe()
}

func (server *Server) Stop() error {
	return server.server.Close()
}
