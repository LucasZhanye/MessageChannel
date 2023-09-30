package core

import (
	"messagechannel/pkg/logger"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// SubScriptionManager
// Subscription include topic and group
// A topic can has many group, a group can has many subscriber
// Only one subscriber of the same subscription group can receive messages
type SubScriptionManager struct {
	subs cmap.ConcurrentMap[string, *topic]
	log  logger.Log
}

type topic struct {
	cmap.ConcurrentMap[string, *group]
}

type group struct {
	cmap.ConcurrentMap[string, *Subscription]
}

type Subscription struct {
	topic  string
	group  string
	client *Client
}

func New(l logger.Log) *SubScriptionManager {
	return &SubScriptionManager{
		subs: cmap.New[*topic](),
		log:  l,
	}
}

func newTopic() *topic {
	return &topic{
		cmap.New[*group](),
	}
}

func newGroup() *group {
	return &group{
		cmap.New[*Subscription](),
	}
}

func NewSubscription(topic, group string, c *Client) *Subscription {
	return &Subscription{
		topic:  topic,
		group:  group,
		client: c,
	}
}

func (m *SubScriptionManager) Add() error {

	return nil

}

func (m *SubScriptionManager) Delete(string) {

}

func (m *SubScriptionManager) Update(string) {

}

func (m *SubScriptionManager) Get(string) {

}

func (m *SubScriptionManager) GetAll() {

}
