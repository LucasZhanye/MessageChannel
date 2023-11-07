package protocol

import (
	"encoding/json"
	"sync"
)

const MESSAGE_RESPONSE_CMDID = "0"

var responsePool = sync.Pool{}

func AcquireResponse(cmdId string) *Response {
	r := responsePool.Get()
	if r == nil {
		return &Response{
			CmdId: cmdId,
		}
	}

	resp := r.(*Response)
	resp.CmdId = cmdId

	return resp
}

func ReleaseResponse(resp *Response) {
	resp.CmdId = ""
	resp.Error = nil
	resp.Connect = nil
	resp.Subscribe = nil
	resp.Unsubscribe = nil
	resp.Publish = nil
	resp.SyncPublishReply = nil
	resp.Ack = nil
	resp.Message = nil
	responsePool.Put(resp)
}

// Response sent from server to client
type Response struct {
	CmdId string
	Error *Error

	Connect          *ConnectResponse
	Subscribe        *SubscribeResponse
	Unsubscribe      *UnsubscribeResponse
	Publish          *PublishResponse
	SyncPublishReply *SyncPublishReplyResponse
	Ack              *AckResponse
	Message          *Message
}

type ConnectResponse struct {
	Identifie string
}

type SubscribeResponse struct {
}

type UnsubscribeResponse struct {
}

type PublishResponse struct {
}

type SyncPublishReplyResponse struct {
}

type AckResponse struct {
}

func (r *Response) SetCmdId(id string) *Response {
	r.CmdId = id
	return r
}

func (r *Response) SetError(err *Error) *Response {
	r.Error = err
	return r
}

func (r *Response) SetConnectResponse(connect *ConnectResponse) *Response {
	r.Connect = connect
	return r
}

func (r *Response) SetSubscribeResponse(sub *SubscribeResponse) *Response {
	r.Subscribe = sub
	return r
}

func (r *Response) SetUnsubscribeResponse(unsub *UnsubscribeResponse) *Response {
	r.Unsubscribe = unsub
	return r
}

func (r *Response) SetPublishResponse(pub *PublishResponse) *Response {
	r.Publish = pub
	return r
}

func (r *Response) SetSyncPublishReplyResponse(spub *SyncPublishReplyResponse) *Response {
	r.SyncPublishReply = spub
	return r
}

func (r *Response) SetAckResponse(ack *AckResponse) *Response {
	r.Ack = ack
	return r
}

func (r *Response) SetMessageResponse(msg *Message) *Response {
	r.Message = msg
	return r
}

func (r *Response) ToJsonString() string {
	jsonByte, _ := json.Marshal(r)
	return string(jsonByte)
}

func (r *Response) ToJsonBytes() []byte {
	ret, _ := json.Marshal(r)

	return ret
}
