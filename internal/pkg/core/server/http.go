package server

import (
	"context"
	"errors"

	"messagechannel/internal/pkg/core"
	v1 "messagechannel/internal/pkg/core/server/api/v1"
	"messagechannel/internal/pkg/core/server/config"
	"messagechannel/internal/pkg/web"
	"messagechannel/pkg/middleware"
	"messagechannel/pkg/safego"

	"messagechannel/internal/pkg/transport/websocket"

	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type HttpServer struct {
	*gin.Engine

	secureServer   *http.Server
	insecureServer *http.Server

	// middlewares []string

	node   *core.Node
	config *config.HttpConfig
}

func NewHttpServer(node *core.Node) *HttpServer {
	return &HttpServer{
		Engine: gin.New(),
		node:   node,
		config: config.NewHttpConfig(),
	}
}

// setup preset for gin
func (s *HttpServer) setup() error {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		s.node.Log.Info("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// check web manager page
	if s.config.Web.Enable {
		if s.config.Web.Password == "" {
			return errors.New("Enable web manager page must be set password on config file")
		}
	}

	return nil
}

// setMiddlewares setting middlewares
func (s *HttpServer) setMiddlewares() {
	s.Engine.Use(gin.Recovery(), gin.Logger())

	s.Engine.Use(middleware.PassParameters("node", s.node))
	s.Engine.Use(middleware.PassParameters("config", s.config))

	if s.config.Web.Enable {
		s.Engine.Use(middleware.ServerStatic(web.FileName, web.Web))
	}
}

func (s *HttpServer) setRoute() {
	e := s.Engine

	apiV1 := e.Group("/api/v1")
	{
		apiV1.GET("/version", v1.VersionHandler)

		// use gin.WrapH to wrap Handler，let handler as gin's HandlerFunc
		apiV1.GET("/conn/websocket", gin.WrapH(websocket.NewWebSocketHandler(s.node)))

		if s.config.Web.Enable {
			apiV1.POST("/login", v1.LoginHandler)

			web := apiV1.Group("/dashboard").Use(middleware.JWT(s.config.Web.Secret))
			{
				web.GET("/sysinfo", v1.SystemInfoHandler)
				web.GET("/clients", v1.ClientInfoHandler)
				web.GET("/subscriptions", v1.SubscriptionHandler)
				web.POST("/subscribe", v1.SubscribeHandler)
				web.POST("/unsubscribe", v1.UnsubscribeHandler)
			}
		}

	}
}

// Initialize Server Initialize
func (s *HttpServer) Initialize() error {
	err := s.setup()
	s.setMiddlewares()
	s.setRoute()

	return err
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
	safego.Execute(s.node.Log, func() {
		for {
			select {
			case <-s.node.GetShutDownChan():
				s.Stop()
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
