package core

import (
	"messagechannel/pkg/protocol"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type ClientInfo struct {
	Address string
	subInfo *subInfo
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

func (cm *ClientInfo) AddSubscription(subscription *protocol.Subscription) error {

	if cm.subInfo == nil {
		cm.subInfo = newSubInfo()
	}

	cm.subInfo.Set(subscription.Topic, subscription)

	return nil
}

func (cm *ClientInfo) RemoveSubscription(subscription *protocol.Subscription) {

}
