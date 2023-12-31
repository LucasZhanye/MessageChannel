package core

import (
	"errors"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"messagechannel/internal/pkg/core/config"
	"messagechannel/pkg/logger"
	"messagechannel/pkg/protocol"
)

// Node represent current MessageChannel
type Node struct {
	Config *config.Config

	Log logger.Log

	engine Engine

	clientManager       *ClientManager
	subscriptionManager *SubscriptionManager
	shutdownChan        chan struct{}
}

// New
func New(log logger.Log) *Node {
	cfg := config.NewConfig()
	if cfg.Name == "" {
		panic("node name must be provided")
	}

	clientManager := NewClientManager(log)
	subscriptionManager := NewSubscriptionManager()

	engine, err := NewEngine(cfg.EngineType, log, subscriptionManager)
	if err != nil {
		panic(err)
	}

	return &Node{
		Config:              cfg,
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

func (n *Node) GetClientManager() *ClientManager {
	return n.clientManager
}

func (n *Node) GetSubscriptionManager() *SubscriptionManager {
	return n.subscriptionManager
}

func (n *Node) Register(client *Client) error {
	return n.clientManager.Add(client)
}

func (n *Node) UnRegister(client *Client) {
	n.clientManager.Remove(client.Identifie)
}

func (n *Node) validateTopicName(topic string) bool {
	rule := n.Config.TopicRule
	// n.Log.Debug("topic rule = %v", rule)

	if rule == "" {
		return true
	}

	matched, err := regexp.MatchString(rule, topic)
	if err != nil {
		n.Log.Error("ValidateTopicName Error = %v", err)
		return false
	}

	return matched
}

func (n *Node) Subscribe(sub *protocol.Subscription) error {
	if !n.validateTopicName(sub.Topic) {
		return errors.New("topic not match rule")
	}

	client, ok := n.clientManager.Get(sub.Identifie)
	if !ok {
		return errors.New("client not register")
	}
	exist := client.Info.CheckSubscription(sub)
	if !exist {
		err := n.engine.Subscribe(sub)
		if err != nil {
			return err
		}
		err = client.Info.AddSubscription(sub)
		if err != nil {
			// exit subscribe goroutinue
			close(sub.ExitChan)
			return err
		}
		err = n.subscriptionManager.Add(sub.Topic, sub.Group, client)
		if err != nil {
			client.Info.RemoveSubscription(sub.Topic)
			return err
		}
	} else {
		return errors.New("Subscription topic exist")
	}

	return nil
}

func (n *Node) Publish(pub *protocol.Publication) error {
	if !n.validateTopicName(pub.Topic) {
		return errors.New("topic not match rule")
	}
	return n.engine.Publish(pub)
}

func (n *Node) SyncPublish(pub *protocol.SyncPublication) (*protocol.Message, error) {
	if !n.validateTopicName(pub.Topic) {
		return nil, errors.New("topic not match rule")
	}

	// Check if the topic has subscribers
	ret := n.subscriptionManager.CheckTopicHasSubscriber(pub.Topic)
	if !ret {
		return nil, errors.New("This topic has no subscribers")
	}

	resp, err := n.engine.SyncPublish(pub)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) SyncPublishReply(reply *protocol.SyncPublicationReply) error {
	return n.engine.SyncPublishReply(reply)
}

func (n *Node) Ack(req *protocol.AckRequest) error {
	return n.engine.Ack(req)
}

func (n *Node) UnSubscribe(unsub *protocol.Unsubscription) error {
	if !n.validateTopicName(unsub.Topic) {
		return errors.New("topic not match rule")
	}
	return n.engine.UnSubscribe(unsub)
}

func (n *Node) Run() {
	err := n.engine.Run()
	if err != nil {
		panic(err)
	}
}
