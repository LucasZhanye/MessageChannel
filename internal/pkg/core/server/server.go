package server

import (
	"messagechannel/internal/pkg/core"
	"messagechannel/internal/pkg/core/server/grpcserver"
	"messagechannel/internal/pkg/core/server/httpserver"
)

type Server struct {
	node *core.Node

	httpServer *httpserver.HttpServer
	grpcServer *grpcserver.GRPCServer
}

func New(node *core.Node) *Server {

	httpServer := httpserver.New(node)

	return &Server{
		node:       node,
		httpServer: httpServer,
	}
}

func (s *Server) preConfig() {

	s.httpServer.Initialize()
}

func (s *Server) Run() {
	s.preConfig()

	s.httpServer.Run()
}
