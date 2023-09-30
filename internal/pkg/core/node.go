package core

import (
	"messagechannel/internal/pkg/engine"
	"os"
	"os/signal"
	"syscall"

	"messagechannel/pkg/logger"
)

// Node represent current MessageChannel
type Node struct {
	Name string

	config Config

	Log logger.Log

	Engine engine.Engine

	SubManager SubScriptionManager

	shutdownChan chan struct{}
}

// New
func NewNode(log logger.Log) *Node {

	return &Node{
		Log:          log,
		shutdownChan: make(chan struct{}),
	}
}

// HandleSignals watch signals to notify exit
func (n *Node) HandleSignals() {
	signalChan := make(chan os.Signal, 1)

	// 监听的信号,ctrl-c:SIGINT,kill pid:SIGTERM
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		n.Log.Info("receive signal: %v, shutdown...", sig)
		close(n.shutdownChan)
		return
	}

}

// GetShutDownChan
func (n *Node) GetShutDownChan() chan struct{} {
	return n.shutdownChan
}
