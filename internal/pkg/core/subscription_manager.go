package core

import (
	"errors"

	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/atomic"
)

type SubscriptionManager struct {
	topics *topic
	total  atomic.Uint32
}

type topic struct {
	cmap.ConcurrentMap[string, *group] // key is topic name
}

type group struct {
	cmap.ConcurrentMap[string, *Subscription] // key is group name
}

type Subscription struct {
	cmap.ConcurrentMap[string, *Client] // key is client's identifie
}

func newgroup() *group {
	return &group{
		cmap.New[*Subscription](),
	}
}

func NewSubscription() *Subscription {
	return &Subscription{
		cmap.New[*Client](),
	}
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		topics: &topic{
			cmap.New[*group](),
		},
	}
}

func (sm *SubscriptionManager) exist(topic, group string) (*Subscription, bool) {
	groups, ok := sm.topics.Get(topic)
	if !ok {
		return nil, false
	}
	subs, ok := groups.Get(group)
	if !ok {
		return nil, false
	}

	return subs, true
}

func (sm *SubscriptionManager) Add(topic, group string, client *Client) error {
	groups, ok := sm.topics.Get(topic)
	if ok {
		clients, ok := groups.Get(group)
		if ok {
			ok := clients.Has(client.Identifie)
			if ok {
				return errors.New("SubscriptionManager add error, because subscription exist.")
			} else {
				// topic,group exist, but client not exist
				clients.Set(client.Identifie, client)
			}
		} else {
			// topic exist， but group not exist, so that client not exist
			g := newgroup()
			sub := NewSubscription()
			sub.Set(client.Identifie, client)
			g.Set(group, sub)
		}
	} else {
		// topic, group, client all not exist
		sub := NewSubscription()
		sub.Set(client.Identifie, client)

		g := newgroup()
		g.Set(group, sub)

		sm.topics.Set(topic, g)

	}

	sm.total.Inc()
	return nil
}

func (sm *SubscriptionManager) Get(topic, group, clientIdentifie string) *Client {
	subscription, ok := sm.exist(topic, group)
	if !ok {
		return nil
	}

	client, ok := subscription.Get(clientIdentifie)
	if !ok {
		return nil
	}

	return client
}

func (sm *SubscriptionManager) GetSubscriptionCount(topic, group string) int {
	subs, ok := sm.exist(topic, group)
	if !ok {
		return 0
	}

	return subs.Count()
}

func (sm *SubscriptionManager) Remove(topic, group, clientIdentifie string) bool {
	subs, ok := sm.exist(topic, group)
	if ok {
		subs.Remove(clientIdentifie)
		sm.total.Dec()
		return true
	}

	return false
}

func (sm *SubscriptionManager) GetAll() *topic {
	return sm.topics
}
