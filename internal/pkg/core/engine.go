package core

import (
	"errors"

	"messagechannel/pkg/logger"
	"messagechannel/pkg/protocol"

	"github.com/spf13/viper"
)

// Engine message engine，support in-memory,rabbitmq,kafka etc.
type Engine interface {
	Run() error

	Close()

	// Publish
	Publish(*protocol.Publication) error

	// Subscribe
	Subscribe(*protocol.Subscription) error

	// UnSubscribe
	UnSubscribe(*protocol.Unsubscription) error

	Ack(*protocol.AckRequest) error
}

func NewEngine(log logger.Log, subscriptionManager *SubscriptionManager) (Engine, error) {
	eType := viper.GetString("engine.type")

	switch eType {
	case "rabbitmq":
		return NewRabbitMqEngine(log, subscriptionManager)
	default:
		return nil, errors.New("Engine type must be provided")
	}
}
