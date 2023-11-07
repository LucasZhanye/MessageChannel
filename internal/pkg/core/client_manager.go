package core

import (
	"fmt"
	"messagechannel/pkg/logger"
	"messagechannel/pkg/protocol"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type ClientManager struct {
	info cmap.ConcurrentMap[string, *Client]
	log  logger.Log
}

func NewClientManager(l logger.Log) *ClientManager {
	return &ClientManager{
		info: cmap.New[*Client](),
		log:  l,
	}
}

func (cm *ClientManager) Add(client *Client) error {

	ok := cm.info.Has(client.Identifie)
	if ok {
		return fmt.Errorf("Client[%s] exist", client.Identifie)
	}

	cm.info.Set(client.Identifie, client)

	return nil
}

func (cm *ClientManager) Remove(identifie string) {
	// first remove client's subscription
	client, ok := cm.Get(identifie)
	if !ok {
		return
	}

	subs := client.Info.SubInfo
	if subs != nil {
		subs.IterCb(func(topic string, subscription *protocol.Subscription) {
			// exit subscribe goroutinue
			close(subscription.ExitChan)
		})
	}

	client.Info.SubInfo = nil
	client.Info = nil

	cm.info.Remove(identifie)
}

func (cm *ClientManager) Get(identifie string) (*Client, bool) {
	return cm.info.Get(identifie)
}

func (cm *ClientManager) GetAll() []*Client {
	clientDtos := make([]*Client, 0)

	cm.info.IterCb(func(key string, v *Client) {
		clientDtos = append(clientDtos, v)
	})

	return clientDtos
}
