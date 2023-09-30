package protocol

import "time"

// Command 客户端和服务端通信命令
type Command struct {
	Id uint32
}

type SubscribeCommand struct {
	Topic string
	Group string
}

type UnsubscribeCommand struct {
	Topic string
}

type PublishCommand struct {
	Topic string
	Data  []byte
}

type PublishWaitCommand struct {
	Topic   string
	Data    []byte
	Timeout time.Duration
}
