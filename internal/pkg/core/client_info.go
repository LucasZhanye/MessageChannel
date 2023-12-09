package core

import (
	"messagechannel/pkg/protocol"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type ClientInfo struct {
	Address     string   `json:"address"`
	ConnectTime int64    `json:"connect_time"`
	SubInfo     *subInfo `json:"subinfo"`
}

type subInfo struct {
	// key means topic name
	cmap.ConcurrentMap[string, *protocol.Subscription]
}

func newSubInfo() *subInfo {
	return &subInfo{
		cmap.New[*protocol.Subscription](),
	}
}

func (ci *ClientInfo) AddSubscription(subscription *protocol.Subscription) error {
	if ci.SubInfo == nil {
		ci.SubInfo = newSubInfo()
	}

	ci.SubInfo.Set(subscription.Topic, subscription)

	return nil
}

func (ci *ClientInfo) CheckSubscription(subscription *protocol.Subscription) bool {
	if ci.SubInfo == nil {
		return false
	}

	return ci.SubInfo.Has(subscription.Topic)
}

func (ci *ClientInfo) GetSubscription(topic string) *protocol.Subscription {
	if ci.SubInfo != nil {
		sub, ok := ci.SubInfo.Get(topic)
		if !ok {
			return nil
		}

		return sub
	}

	return nil
}

func (ci *ClientInfo) RemoveSubscription(topic string) {
	if ci.SubInfo != nil {
		subscription := ci.GetSubscription(topic)

		if subscription != nil {
			ci.SubInfo.Remove(subscription.Topic)
			close(subscription.ExitChan)
		}
	}
}

func (ci *ClientInfo) GetAllSubscription() map[string]*protocol.Subscription {
	if ci.SubInfo != nil {
		return ci.SubInfo.Items()
	}

	return nil
}
