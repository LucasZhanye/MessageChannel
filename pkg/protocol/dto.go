package protocol

import (
	"strings"
	"time"
)

type Publication struct {
	Identifie string
	Topic     string
	Data      []byte
	MetaData  map[string]string
}

type SyncPublication struct {
	Identifie string
	Topic     string
	Data      []byte
	Timeout   time.Duration
}

type SyncPublicationReply struct {
	Identifie string
	MessageId string
	Data      []byte
	Header    map[string]any
}

type Subscription struct {
	Identifie string
	Topic     string
	Group     string
	ExitChan  chan struct{}
}

type Unsubscription struct {
	Identifie string
	Topic     string
	Group     string
}

func NewSubscription(identifie string, topic string, group string) *Subscription {
	return &Subscription{
		Identifie: identifie,
		Topic:     strings.ToLower(topic),
		Group:     strings.ToLower(group),
		ExitChan:  make(chan struct{}),
	}
}

func NewUnSubscription(identifie, topic, group string) *Unsubscription {
	return &Unsubscription{
		Identifie: identifie,
		Topic:     strings.ToLower(topic),
		Group:     strings.ToLower(group),
	}
}
