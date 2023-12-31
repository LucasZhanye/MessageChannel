package protocol

import (
	"encoding/json"
	"time"
)

type AckResult uint8

const (
	ACKRESULT_ACK AckResult = iota + 1
	ACKRESULT_NACK
	ACKRESULT_REJECT
)

// Request sent from client to server
type Request struct {
	CmdId            string
	Connect          *ConnectRequest
	Subscribe        *SubscribeRequest
	Unsubscribe      *UnsubscribeRequest
	Publish          *PublishRequest
	SyncPublish      *SyncPublishRequest
	SyncPublishReply *SyncPublishReplyRequest
	Ack              *AckRequest
}

type ConnectRequest struct {
	Name string `json:"name"`
}

type SubscribeRequest struct {
	Identifie string `json:"identifie"`
	Topic     string `json:"topic"`
	Group     string `json:"group"`
}

type UnsubscribeRequest struct {
	Topic string `json:"topic"`
	Group string `json:"group"`
}

type PublishRequest struct {
	Topic string `json:"topic"`
	Data  []byte `json:"data"`
}

type SyncPublishRequest struct {
	Topic    string        `json:"topic"`
	Data     []byte        `json:"data"`
	Timeouut time.Duration `json:"timeout"`
}

type SyncPublishReplyRequest struct {
	MessageId string `json:"messageId"`
	Data      []byte `json:"data"`
	Header map[string]any `json:"header"`
}

type AckRequest struct {
	MessageId string    `json:"messageId"`
	Result    AckResult `json:"result"`
}

func ConvertToRequest(msg []byte) *Request {
	req := &Request{}
	err := json.Unmarshal(msg, req)
	if err != nil {
		return nil
	}

	return req
}
