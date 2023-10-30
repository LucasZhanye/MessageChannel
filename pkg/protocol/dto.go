package protocol

import "strings"

type Publication struct {
	Identifie string
	Topic     string
	Data      []byte
}

type Subscription struct {
	Identifie string
	Topic     string
	Group     string
	ExitChan  chan struct{}
}

type Unsubscription struct {
	Topic string
	Group string
}

func NewSubscription(identifie string, topic string, group string) *Subscription {
	return &Subscription{
		Identifie: identifie,
		Topic:     strings.ToLower(topic),
		Group:     strings.ToLower(group),
		ExitChan:  make(chan struct{}),
	}
}
