package core

import (
	"messagechannel/internal/pkg/transport"
	"messagechannel/pkg/protocol"
	"strings"
	"sync"
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
	Identifie string      `json:"identifie"`
	Info      *ClientInfo `json:"info"`

	transport transport.Transport

	status status

	writeLock sync.Mutex
}

func NewClient(node *Node, info *ClientInfo, transport transport.Transport) *Client {
	client := &Client{
		node:      node,
		Info:      info,
		transport: transport,
		status:    CLIENT_STATUS_CONNECTING,
	}

	return client
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
		c.node.Log.Error("[client-%s]message can not convert to request struct", c.Info.Address)
		return transport.CloseEventBadRequest, false
	}

	var err error
	if req.Connect != nil {
		err = c.handleConnect(req)
	} else if req.Publish != nil {
		err = c.handlePublish(req)
	} else if req.SyncPublish != nil {
		go c.handleSyncPublish(req)
	} else if req.SyncPublishReply != nil {
		err = c.handleSyncPublishReply(req)
	} else if req.Subscribe != nil {
		err = c.handleSubscribe(req)
	} else if req.Unsubscribe != nil {
		err = c.handleUnsubscribe(req)
	} else if req.Ack != nil {
		err = c.handleAck(req)
	}

	if err != nil {
		c.node.Log.Error("[client-%s]handle request (%+v) error: %v", c.Info.Address, req, err)
		return transport.CloseEventInternalError, false
	}

	return nil, true
}

func (c *Client) handleConnect(req *protocol.Request) error {
	connectReq := req.Connect

	response := protocol.AcquireResponse(req.CmdId)
	defer protocol.ReleaseResponse(response)

	if connectReq.Name == "" {
		return c.Send(response.SetError(protocol.ErrConnect.WithReason("client name is null")).ToJsonBytes())
	}

	c.Identifie = connectReq.Name

	// Add to ClientManager
	err := c.node.Register(c)
	if err != nil {
		return c.Send(response.SetError(protocol.ErrConnect.WithReason(err.Error())).ToJsonBytes())
	}

	c.status = CLIENT_STATUS_CONNECTED

	resp := response.SetConnectResponse(&protocol.ConnectResponse{
		Identifie: c.Identifie,
	}).ToJsonBytes()

	err = c.Send(resp)
	if err != nil {
		return err
	}

	c.node.Log.Info("client[%s] connect success!", c.Identifie)
	return nil
}

func (c *Client) handleSubscribe(req *protocol.Request) error {
	var err error
	var subscription *protocol.Subscription

	subReq := req.Subscribe

	c.node.Log.Debug("subReq = %+v", subReq)

	response := protocol.AcquireResponse(req.CmdId)
	defer protocol.ReleaseResponse(response)

	if subReq.Identifie != c.Identifie {
		response.SetError(protocol.ErrClientIdMatch)
		goto Reply
	}

	if subReq.Topic == "" {
		response.SetError(protocol.ErrTopicIsNull)
		goto Reply
	}
	if subReq.Group == "" {
		response.SetError(protocol.ErrGroupIsNull)
		goto Reply
	}

	subscription = protocol.NewSubscription(subReq.Identifie, subReq.Topic, subReq.Group)
	err = c.node.Subscribe(subscription)
	if err != nil {
		response.SetError(protocol.ErrSubscribe.WithReason(err.Error()))
		goto Reply
	}

	response.SetSubscribeResponse(&protocol.SubscribeResponse{})
Reply:
	return c.Send(response.ToJsonBytes())
}

func (c *Client) handleUnsubscribe(req *protocol.Request) error {
	var err error
	var unsub *protocol.Unsubscription

	unsubReq := req.Unsubscribe

	c.node.Log.Debug("unsubReq = %+v", unsubReq)

	response := protocol.AcquireResponse(req.CmdId)
	defer protocol.ReleaseResponse(response)

	if unsubReq.Topic == "" {
		response.SetError(protocol.ErrTopicIsNull)
		goto Reply
	}
	if unsubReq.Group == "" {
		response.SetError(protocol.ErrGroupIsNull)
		goto Reply
	}

	unsub = protocol.NewUnSubscription(c.Identifie, unsubReq.Topic, unsubReq.Group)

	err = c.node.UnSubscribe(unsub)
	if err != nil {
		response.SetError(protocol.ErrUnsubscribe.WithReason(err.Error()))
		goto Reply
	}

	response.SetUnsubscribeResponse(&protocol.UnsubscribeResponse{})
Reply:
	return c.Send(response.ToJsonBytes())
}

func (c *Client) handlePublish(req *protocol.Request) error {
	var err error
	var publication *protocol.Publication

	publishReq := req.Publish

	c.node.Log.Debug("publishReq = %+v", publishReq)

	response := protocol.AcquireResponse(req.CmdId)
	defer protocol.ReleaseResponse(response)

	if publishReq.Topic == "" {
		response.SetError(protocol.ErrTopicIsNull)
		goto Reply
	}
	if len(publishReq.Data) == 0 {
		response.SetError(protocol.ErrDataIsNull)
		goto Reply
	}

	publication = &protocol.Publication{
		Identifie: c.Identifie,
		Topic:     strings.ToLower(publishReq.Topic),
		Data:      publishReq.Data,
	}

	err = c.node.Publish(publication)
	if err != nil {
		response.SetError(protocol.ErrPublish.WithReason(err.Error()))
		goto Reply
	}

Reply:
	return c.Send(response.ToJsonBytes())
}

func (c *Client) handleSyncPublish(req *protocol.Request) error {
	var err error
	var message *protocol.Message
	var syncPublication *protocol.SyncPublication

	publishReq := req.SyncPublish

	c.node.Log.Debug("sync publishReq = %+v", publishReq)

	response := protocol.AcquireResponse(req.CmdId)
	defer protocol.ReleaseResponse(response)

	if publishReq.Topic == "" {
		response.SetError(protocol.ErrTopicIsNull)
		goto Reply
	}
	if len(publishReq.Data) == 0 {
		response.SetError(protocol.ErrDataIsNull)
		goto Reply
	}

	syncPublication = &protocol.SyncPublication{
		Identifie: c.Identifie,
		Topic:     strings.ToLower(publishReq.Topic),
		Data:      publishReq.Data,
		Timeout:   publishReq.Timeouut,
	}

	message, err = c.node.SyncPublish(syncPublication)
	if err != nil {
		response.SetError(protocol.ErrPublish.WithReason(err.Error()))
		goto Reply
	}

	response.SetMessageResponse(message)
Reply:
	return c.Send(response.ToJsonBytes())

}

func (c *Client) handleSyncPublishReply(req *protocol.Request) error {
	var err error
	var reply *protocol.SyncPublicationReply
	replyReq := req.SyncPublishReply

	c.node.Log.Debug("replyReq = %+v", replyReq)

	response := protocol.AcquireResponse(req.CmdId)
	defer protocol.ReleaseResponse(response)

	if replyReq.MessageId == "" {
		response.SetError(protocol.ErrSyncPublishReply.WithReason("message id is null"))
		goto Reply
	}

	reply = &protocol.SyncPublicationReply{
		Identifie: c.Identifie,
		MessageId: replyReq.MessageId,
		Data:      replyReq.Data,
		Header:    replyReq.Header,
	}

	err = c.node.SyncPublishReply(reply)
	if err != nil {
		response.SetError(protocol.ErrSyncPublishReply.WithReason(err.Error()))
		goto Reply
	}

	response.SetSyncPublishReplyResponse(&protocol.SyncPublishReplyResponse{})
Reply:
	return c.Send(response.ToJsonBytes())
}

func (c *Client) handleAck(req *protocol.Request) error {
	var err error

	ackReq := req.Ack

	c.node.Log.Debug("ackReq = %+v", ackReq)

	response := protocol.AcquireResponse(req.CmdId)
	defer protocol.ReleaseResponse(response)

	if ackReq.MessageId == "" {
		response.SetError(protocol.ErrAck.WithReason("message id is null"))
		goto Reply
	}

	err = c.node.Ack(ackReq)
	if err != nil {
		response.SetError(protocol.ErrAck.WithReason(err.Error()))
		goto Reply
	}

	response.SetAckResponse(&protocol.AckResponse{})
Reply:
	return c.Send(response.ToJsonBytes())
}
