package server

import (
	"messagechannel/internal/pkg/core"
)

type Server struct {
	node *core.Node

	httpServer *HttpServer
	grpcServer *GRPCServer
}

func New(node *core.Node) *Server {

	httpServer := NewHttpServer(node)

	return &Server{
		node:       node,
		httpServer: httpServer,
	}
}

func (s *Server) preConfig() error {
	return s.httpServer.Initialize()
}

func (s *Server) Run() error {
	err := s.preConfig()
	if err != nil {
		return err
	}
	return s.httpServer.Run()
}
