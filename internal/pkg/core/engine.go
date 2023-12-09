package core

import (
	"errors"

	"messagechannel/pkg/logger"
	"messagechannel/pkg/protocol"
)

// Engine message engine，support in-memory,rabbitmq,kafka etc.
type Engine interface {
	Run() error

	Close()

	// Publish
	Publish(*protocol.Publication) error

	// SyncPublish
	SyncPublish(*protocol.SyncPublication) (*protocol.Message, error)

	// SyncPublishReply
	SyncPublishReply(*protocol.SyncPublicationReply) error

	// Subscribe
	Subscribe(*protocol.Subscription) error

	// UnSubscribe
	UnSubscribe(*protocol.Unsubscription) error

	Ack(*protocol.AckRequest) error
}

func NewEngine(engineType string, log logger.Log, subscriptionManager *SubscriptionManager) (Engine, error) {
	switch engineType {
	case "rabbitmq":
		return NewRabbitMqEngine(log, subscriptionManager)
	default:
		return nil, errors.New("Engine type must be provided")
	}
}
