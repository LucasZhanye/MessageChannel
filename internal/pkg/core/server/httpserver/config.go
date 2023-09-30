package httpserver

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HttpAddress       string
	HttpPort          string
	ShutdownDelayTime time.Duration // wait for shutdown
}

func setDefaultConfig(prefix string) {

	defaultConfig := map[string]any{
		prefix + ".insecure_address":    "127.0.0.1",
		prefix + ".insecure_port":       "7788",
		prefix + ".shutdown_delay_time": 5 * time.Second,
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}

func NewConfig(module string) *Config {

	// setting default config
	setDefaultConfig(module)

	// get config from viper
	httpAddress := viper.GetString(module + ".insecure_address")
	httpPort := viper.GetString(module + ".insecure_port")
	shutdownDelayTime := viper.GetDuration(module + ".shutdown_delay_time")

	return &Config{
		HttpAddress:       httpAddress,
		HttpPort:          httpPort,
		ShutdownDelayTime: shutdownDelayTime,
	}
}
