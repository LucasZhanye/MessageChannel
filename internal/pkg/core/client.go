package core

import (
	"messagechannel/internal/pkg/transport"
	"messagechannel/pkg/protocol"
	"strings"
	"sync"

	"github.com/lithammer/shortuuid/v4"
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

	// generate by client name and uuid, such as clientA-xxxxx
	Identifie string

	info      *ClientInfo
	transport transport.Transport

	status status

	writeLock sync.Mutex
}

func NewClient(node *Node, info *ClientInfo, transport transport.Transport) *Client {
	client := &Client{
		node:      node,
		info:      info,
		transport: transport,
		status:    CLIENT_STATUS_CONNECTING,
	}

	return client
}

func (c *Client) genClientId() string {
	return shortuuid.New()
}

func (c *Client) setIdentifie(name, cid string) string {
	return name + "-" + cid
}

func (c *Client) GetIdentifie() string {
	return c.Identifie
}

// Send send data to client
func (c *Client) Send(data []byte) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	return c.transport.Write(data)
}

// CLose
func (c *Client) Close(event *transport.CloseEvent) error {
	if c.status == CLIENT_STATUS_CONNECTED {
		c.node.Log.Info("Closing Client[%s]", c.Identifie)

		// remove from client manager
		c.node.UnRegister(c)

		c.status = CLIENT_STATUS_CLOSED
		// 关闭连接
		err := c.transport.Close(event)
		if err != nil {
			return err
		}
		c.node.Log.Info("Closed Client[%s]", c.Identifie)
	}

	return nil
}

// HandleRequest handle request from client
// such as subscribe,unsubscribe,publish etc.
func (c *Client) HandleRequest(msg []byte) (*transport.CloseEvent, bool) {
	if c.status == CLIENT_STATUS_CLOSED {
		return transport.CloseEventConnectionClosed, false
	}

	// convert to request struct
	req := protocol.ConvertToRequest(msg)
	if req == nil {
		c.node.Log.Error("[client-%s]message can not convert to request struct", c.info.Address)
		return transport.CloseEventBadRequest, false
	}

	var err error
	if req.Connect != nil {
		err = c.handleConnect(req)
	} else if req.Publish != nil {
		err = c.handlePublish(req)
	} else if req.PublishWait != nil {
		err = c.handlePublishWait(req)
	} else if req.Subscribe != nil {
		err = c.handleSubscribe(req)
	} else if req.Unsubscribe != nil {
		err = c.handleUnsubscribe(req)
	} else if req.Ack != nil {
		err = c.handleAck(req)
	}

	if err != nil {
		c.node.Log.Error("[client-%s]handle request (%+v) error: %v", c.info.Address, req, err)
		return transport.CloseEventInternalError, false
	}

	return nil, true
}

func (c *Client) handleConnect(req *protocol.Request) error {
	connectReq := req.Connect

	var reply *protocol.Reply

	var clientName string = c.info.Address

	if connectReq.Name != "" {
		clientName = connectReq.Name
	}

	cid := c.genClientId()
	c.Identifie = c.setIdentifie(clientName, cid)

	reply = protocol.ReplyPool.GetResponseReply(protocol.OK.WithType(protocol.RESPONSE_CONNECT).WithMetaData(map[string]string{"identifie": c.Identifie}))
	defer protocol.ReplyPool.ReleaseResponseReply(reply)

	// Add to ClientManager
	err := c.node.Register(c)
	if err != nil {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrParam.WithType(protocol.RESPONSE_CONNECT).WithMessage(err.Error()))
		return c.Send(reply.SetId(req.Id).ToJsonBytes())
	}

	c.status = CLIENT_STATUS_CONNECTED

	err = c.Send(reply.SetId(req.Id).ToJsonBytes())
	if err != nil {
		return err
	}

	c.node.Log.Info("client[%s] connect success!", c.Identifie)
	return nil
}

func (c *Client) handleSubscribe(req *protocol.Request) error {
	var reply *protocol.Reply
	var err error
	var subscription *protocol.Subscription

	subReq := req.Subscribe

	c.node.Log.Debug("subReq = %+v", subReq)

	reply = protocol.ReplyPool.GetResponseReply(protocol.OK.WithType(protocol.RESPONSE_SUBSCRIBE))
	defer protocol.ReplyPool.ReleaseResponseReply(reply)

	if subReq.Identifie != c.Identifie {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrParam.WithType(protocol.RESPONSE_SUBSCRIBE).WithMessage("client id not match"))
		goto Reply
	}
	if subReq.Topic == "" {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrParam.WithType(protocol.RESPONSE_SUBSCRIBE).WithMessage("topic is nul"))
		goto Reply
	}
	if subReq.Group == "" {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrParam.WithType(protocol.RESPONSE_SUBSCRIBE).WithMessage("group is nul"))
		goto Reply
	}

	subscription = protocol.NewSubscription(subReq.Identifie, subReq.Topic, subReq.Group)
	err = c.node.Subscribe(subscription)
	if err != nil {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrServer.WithType(protocol.RESPONSE_SUBSCRIBE).WithMessage(err.Error()))
		goto Reply
	}

Reply:
	return c.Send(reply.SetId(req.Id).ToJsonBytes())
}

func (c *Client) handleUnsubscribe(req *protocol.Request) error {
	var reply *protocol.Reply
	var err error
	var unsub *protocol.Unsubscription

	unsubReq := req.Unsubscribe

	c.node.Log.Debug("unsubReq = %+v", unsubReq)

	reply = protocol.ReplyPool.GetResponseReply(protocol.OK.WithType(protocol.RESPONSE_UNSUBSCRIBE))
	defer protocol.ReplyPool.ReleaseResponseReply(reply)

	if unsubReq.Topic == "" {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrParam.WithType(protocol.RESPONSE_UNSUBSCRIBE).WithMessage("topic is nul"))
		goto Reply
	}
	if unsubReq.Group == "" {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrParam.WithType(protocol.RESPONSE_UNSUBSCRIBE).WithMessage("group is nul"))
		goto Reply
	}

	unsub = protocol.NewUnSubscription(unsubReq.Topic, unsubReq.Group)
	err = c.node.UnSubscribe(unsub)
	if err != nil {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrServer.WithType(protocol.RESPONSE_UNSUBSCRIBE).WithMessage(err.Error()))
		goto Reply
	}

Reply:
	return c.Send(reply.SetId(req.Id).ToJsonBytes())
}

func (c *Client) handlePublish(req *protocol.Request) error {
	var reply *protocol.Reply
	var err error
	var publication *protocol.Publication

	publishReq := req.Publish

	c.node.Log.Debug("publishReq = %+v", publishReq)

	reply = protocol.ReplyPool.GetResponseReply(protocol.OK.WithType(protocol.RESPONSE_PUBLISH))
	defer protocol.ReplyPool.ReleaseResponseReply(reply)

	if publishReq.Topic == "" || len(publishReq.Data) == 0 {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrParam.WithType(protocol.RESPONSE_PUBLISH).WithMessage("topic or data is nul"))
		goto Reply
	}

	publication = &protocol.Publication{
		Identifie: c.Identifie,
		Topic:     strings.ToLower(publishReq.Topic),
		Data:      publishReq.Data,
	}

	err = c.node.Publish(publication)
	if err != nil {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrServer.WithType(protocol.RESPONSE_PUBLISH).WithMessage(err.Error()))
		goto Reply
	}

Reply:
	return c.Send(reply.SetId(req.Id).ToJsonBytes())
}

func (c *Client) handlePublishWait(req *protocol.Request) error {
	return nil
}

func (c *Client) handleAck(req *protocol.Request) error {
	var reply *protocol.Reply
	var err error

	ackReq := req.Ack

	c.node.Log.Debug("ackReq = %+v", ackReq)

	reply = protocol.ReplyPool.GetResponseReply(protocol.OK.WithType(protocol.RESPONSE_ACK))
	defer protocol.ReplyPool.ReleaseResponseReply(reply)

	if ackReq.MessageId == "" {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrParam.WithType(protocol.RESPONSE_ACK).WithMessage("message id is nul"))
		goto Reply
	}

	err = c.node.Ack(ackReq)
	if err != nil {
		reply = protocol.ReplyPool.GetResponseReply(protocol.ErrServer.WithType(protocol.RESPONSE_ACK).WithMessage(err.Error()))
		goto Reply
	}

Reply:
	return c.Send(reply.SetId(req.Id).ToJsonBytes())
}
