package core

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"messagechannel/pkg/logger"
	"messagechannel/pkg/protocol"

	"github.com/spf13/viper"
)

// Node represent current MessageChannel
type Node struct {
	Name string

	Log logger.Log

	engine Engine

	clientManager       *ClientManager
	subscriptionManager *SubscriptionManager
	shutdownChan        chan struct{}
}

// New
func New(log logger.Log) *Node {
	name := viper.GetString("name")
	if name == "" {
		panic("node name must be provided")
	}

	clientManager := NewClientManager(log)
	subscriptionManager := NewSubscriptionManager()

	engine, err := NewEngine(log, subscriptionManager)
	if err != nil {
		panic(err)
	}

	return &Node{
		Name:                name,
		Log:                 log,
		engine:              engine,
		clientManager:       clientManager,
		subscriptionManager: subscriptionManager,
		shutdownChan:        make(chan struct{}),
	}
}

// HandleSignals watch signals to notify exit
func (n *Node) HandleSignals() {
	signalChan := make(chan os.Signal, 1)

	// 监听的信号,ctrl-c:SIGINT,kill pid:SIGTERM, for windows :os.Intterrupt
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case sig := <-signalChan:
		n.Log.Info("receive signal: %v, shutdown...", sig)

		n.engine.Close()

		close(n.shutdownChan)

		return
	}
}

// GetShutDownChan
func (n *Node) GetShutDownChan() chan struct{} {
	return n.shutdownChan
}

func (n *Node) Register(client *Client) error {
	return n.clientManager.Add(client)
}

func (n *Node) UnRegister(client *Client) {
	n.clientManager.Remove(client.Identifie)
}

func (n *Node) Subscribe(sub *protocol.Subscription) error {

	err := n.engine.Subscribe(sub)
	if err != nil {
		return err
	}

	client, ok := n.clientManager.Get(sub.Identifie)
	if !ok {
		// exit subscribe goroutinue
		close(sub.ExitChan)

		return errors.New("client not register")
	}

	err = client.info.AddSubscription(sub)
	if err != nil {
		// exit subscribe goroutinue
		close(sub.ExitChan)
		return err
	}

	err = n.subscriptionManager.Add(sub.Topic, sub.Group, client)
	if err != nil {
		// exit subscribe goroutinue
		close(sub.ExitChan)
		client.info.RemoveSubscription(sub)
		return err
	}
	return nil
}

func (n *Node) Publish(pub *protocol.Publication) error {
	return n.engine.Publish(pub)
}

func (n *Node) Ack(req *protocol.AckRequest) error {
	return n.engine.Ack(req)
}

func (n *Node) UnSubscribe(unsub *protocol.Unsubscription) error {

	return nil
}

func (n *Node) Run() {
	err := n.engine.Run()
	if err != nil {
		panic(err)
	}
}
