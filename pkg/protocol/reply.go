package protocol

import (
	"encoding/json"
	"sync"
)

type ReplyType uint8

const (
	REPLY_MESSAGE ReplyType = iota + 1
	REPLY_RESPONSE
)

// Reply sent from server to client,include message or response
type Reply struct {
	Id string `json:"id"`

	Type ReplyType `json:"type"`

	Response *Response `json:"response"`

	Message *Message `json:"message"`
}

type pool struct {
	MessagePool sync.Pool

	ResponsePool sync.Pool
}

var ReplyPool = &pool{}

func (p *pool) GetResponseReply(response *Response) *Reply {
	r := p.ResponsePool.Get()
	if r == nil {
		return &Reply{
			Type:     REPLY_RESPONSE,
			Response: response,
		}
	}
	reply := r.(*Reply)
	reply.Response = response

	return reply
}

func (p *pool) ReleaseResponseReply(r *Reply) {
	r.Response = nil
	p.ResponsePool.Put(r)
}

func (p *pool) GetMessageReply(message *Message) *Reply {
	r := p.MessagePool.Get()
	if r == nil {
		return &Reply{
			Type:    REPLY_MESSAGE,
			Message: message,
		}
	}
	reply := r.(*Reply)
	reply.Message = message

	return reply
}

func (p *pool) ReleaseMessageReply(r *Reply) {
	r.Message = nil
	p.MessagePool.Put(r)
}

func (r *Reply) SetId(id string) *Reply {
	r.Id = id
	return r
}

func (r *Reply) ToJsonString() string {
	jsonByte, _ := json.Marshal(r)
	return string(jsonByte)
}

func (r *Reply) ToJsonBytes() []byte {
	ret, _ := json.Marshal(r)

	return ret
}
