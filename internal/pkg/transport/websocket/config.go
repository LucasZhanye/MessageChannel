package websocket

import (
	"messagechannel/internal/pkg/transport"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const moduleName = "transport.websocket"

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
func setDefaultConfig() {

	defaultConfig := map[string]any{
		moduleName + ".read_buffer_size":      1024,
		moduleName + ".write_buffer_size":     1024,
		moduleName + ".use_write_buffer_pool": false,
		moduleName + ".enable_compression":    false,
		moduleName + ".compression_level":     1,
		moduleName + ".start_ping":            false,
		moduleName + ".ping_interval":         time.Duration(30 * time.Second),
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}

// NewWebSocketConfig
func NewWebSocketConfig() *WebSocketConfig {

	setDefaultConfig()

	readBufferSize := viper.GetInt(moduleName + ".read_buffer_size")
	writeBufferSize := viper.GetInt(moduleName + ".write_buffer_size")
	useWriteBufferPool := viper.GetBool(moduleName + ".use_write_buffer_pool")
	enableCompression := viper.GetBool(moduleName + ".enable_compression")
	compressionLevel := viper.GetInt(moduleName + ".compression_level")
	startPing := viper.GetBool(moduleName + ".start_ping")
	pingInterval := viper.GetDuration(moduleName + ".ping_interval")

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
