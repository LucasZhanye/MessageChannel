package websocket

import (
	"messagechannel/internal/pkg/transport"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// WebSocketConfig
type WebSocketConfig struct {
	// size of read buffer
	ReadBufferSize int

	// size of write buffer
	WriteBufferSize int

	// whether to use the write buffer pool
	UseWriteBufferPool bool

	// whether compression is supported
	EnableCompression bool

	// compression level, See the compress/flate package for a description ofcompression levels.
	CompressionLevel int

	// whether start ping
	StartPing bool

	// interval of ping
	PingInterval time.Duration

	// Cross-domain detection
	CheckOrigin func(r *http.Request) bool
}

// setDefaultConfig
func setDefaultConfig(prefix string) {

	defaultConfig := map[string]any{
		prefix + ".read_buffer_size":      1024,
		prefix + ".write_buffer_size":     1024,
		prefix + ".use_write_buffer_pool": false,
		prefix + ".enable_compression":    false,
		prefix + ".compression_level":     1,
		prefix + ".start_ping":            false,
		prefix + ".ping_interval":         time.Duration(30 * time.Second),
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}

// NewWebSocketConfig
func NewWebSocketConfig(module string) *WebSocketConfig {

	setDefaultConfig(module)

	readBufferSize := viper.GetInt(module + ".read_buffer_size")
	writeBufferSize := viper.GetInt(module + ".write_buffer_size")
	useWriteBufferPool := viper.GetBool(module + ".use_write_buffer_pool")
	enableCompression := viper.GetBool(module + ".enable_compression")
	compressionLevel := viper.GetInt(module + ".compression_level")
	startPing := viper.GetBool(module + ".start_ping")
	pingInterval := viper.GetDuration(module + ".ping_interval")

	return &WebSocketConfig{
		ReadBufferSize:     readBufferSize,
		WriteBufferSize:    writeBufferSize,
		UseWriteBufferPool: useWriteBufferPool,
		EnableCompression:  enableCompression,
		CheckOrigin:        transport.CheckOriginFunc(),
		CompressionLevel:   compressionLevel,
		StartPing:          startPing,
		PingInterval:       pingInterval,
	}
}
