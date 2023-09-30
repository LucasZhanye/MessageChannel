package engine

// Engine 消息中间件引擎，支持Memory
type Engine interface {
	// Publish 发布消息
	Publish()

	// Subscribe 订阅主题
	Subscribe()
}
