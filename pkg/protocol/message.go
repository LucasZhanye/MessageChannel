package protocol

import "github.com/lithammer/shortuuid/v4"

// Message push by server to client
type Message struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Payload   []byte `json:"payload"`
}

func NewMessage(id string, timestamp int64, payload []byte) *Message {
	return &Message{
		Id:        id,
		Timestamp: timestamp,
		Payload:   payload,
	}
}

func GenerateMessageId() string {
	return shortuuid.New()
}
