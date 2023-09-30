package core

import (
	"fmt"
	"messagechannel/internal/pkg/transport"
)

type status uint8

const (
	CLIENT_STATUS_CONNECTING status = iota + 1
	CLIENT_STATUS_CONNECTED
	CLIENT_STATUS_CLOSED
)

// Client
type Client struct {
	node *Node

	transport transport.Transport

	subs []*Subscription

	status status
}

func NewClient(node *Node, transport transport.Transport) *Client {
	return &Client{
		node:      node,
		transport: transport,
		status:    CLIENT_STATUS_CONNECTING,
	}
}

// Send send data to client
func (c *Client) Send(data []byte) error {
	return c.transport.Write(data)
}

// CLose
func (c *Client) Close(event transport.CloseEvent) error {
	if c.status == CLIENT_STATUS_CONNECTED {
		// 取消订阅
		for _, sub := range c.subs {
			// 删除订阅信息
			c.node.SubManager.Delete("")
			// TODO: 取消engine的订阅
			fmt.Println(sub)
		}
		c.status = CLIENT_STATUS_CLOSED
		// 关闭连接
		return c.transport.Close(event)
	}

	return nil
}

// HandleEvent handle event from other client
// such as subscribe,unsubscribe,publish etc.
func (c *Client) HandleEvent() {

}

func (c *Client) handleSubscribe() {

}

func (c *Client) handleUnsubscribe() {

}

func (c *Client) handlePublish() {

}

func (c *Client) handlePublishWait() {

}
