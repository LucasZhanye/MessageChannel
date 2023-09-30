package httpserver

import (
	"context"
	"errors"

	"messagechannel/internal/pkg/core"
	v1 "messagechannel/internal/pkg/core/server/api/v1"

	"messagechannel/internal/pkg/transport/websocket"

	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var moduleName string = "httpserver"

type HttpServer struct {
	*gin.Engine

	secureServer   *http.Server
	insecureServer *http.Server

	middlewares []string

	node   *core.Node
	config *Config
}

func New(node *core.Node) *HttpServer {

	return &HttpServer{
		Engine: gin.New(),
		node:   node,
		config: NewConfig(moduleName),
	}
}

// setup preset for gin
func (s *HttpServer) setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		s.node.Log.Info("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// setMiddlewares setting middlewares
func (s *HttpServer) setMiddlewares() {

}

func (s *HttpServer) setRoute() {
	e := s.Engine
	apiV1 := e.Group("/api/v1")
	{
		apiV1.GET("/version", v1.VersionHandler)

		// use gin.WrapH to wrap Handler，let handler as gin's HandlerFunc
		apiV1.GET("/conn/websocket", gin.WrapH(websocket.NewWebSocketHandler(s.node)))
	}
}

// Initialize Server Initialize
func (s *HttpServer) Initialize() {
	s.setup()
	s.setMiddlewares()
	s.setRoute()
}

// Run start server
func (s *HttpServer) Run() error {
	s.insecureServer = &http.Server{
		Addr:    net.JoinHostPort(s.config.HttpAddress, s.config.HttpPort),
		Handler: s,
	}

	// TODO: secureServer handle

	var eg errgroup.Group

	eg.Go(func() error {
		s.node.Log.Info("Start to listening on http address: %s", s.insecureServer.Addr)

		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.node.Log.Fatal(err.Error())

			return err
		}
		return nil
	})

	// TODO: start TLS Server

	// notify node exit
	eg.Go(func() error {
		for {
			select {
			case <-s.node.GetShutDownChan():
				s.Stop()
				return nil
			}
		}
	})

	if err := eg.Wait(); err != nil {
		s.node.Log.Fatal(err.Error())
	}

	return nil
}

// Stop graceful shutdown server
func (s *HttpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownDelayTime)
	defer cancel()

	s.node.Log.Info("Start shutdown http server...")
	if err := s.insecureServer.Shutdown(ctx); err != nil {
		s.node.Log.Warn("Shutdown insecure server failed: %s", err.Error())
	}

	// TODO: shutdown secure server

	s.node.Log.Info("Shutdown http server finish!")
}
