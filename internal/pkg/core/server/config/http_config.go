package config

import (
	"time"

	"github.com/spf13/viper"
)

const moduleName string = "server.http"

type HttpConfig struct {
	HttpAddress       string
	HttpPort          string
	ShutdownDelayTime time.Duration // wait for shutdown
	Web               *webConfig
}

type webConfig struct {
	Enable     bool
	Password   string
	Secret     string
	ExpireTime time.Duration
}

func setDefaultConfig() {
	defaultConfig := map[string]any{
		moduleName + ".insecure_address":    "127.0.0.1",
		moduleName + ".insecure_port":       "7788",
		moduleName + ".shutdown_delay_time": 5 * time.Second,
		moduleName + ".web.enable":          false,
		moduleName + ".web.secret":          "messagechannel",
		moduleName + ".web.expire_time":     30 * time.Minute,
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}

func NewHttpConfig() *HttpConfig {

	// setting default config
	setDefaultConfig()

	// get config from viper
	httpAddress := viper.GetString(moduleName + ".insecure_address")
	httpPort := viper.GetString(moduleName + ".insecure_port")
	shutdownDelayTime := viper.GetDuration(moduleName + ".shutdown_delay_time")
	webEnable := viper.GetBool(moduleName + ".web.enable")
	webPassword := viper.GetString(moduleName + ".web.password")
	webSecret := viper.GetString(moduleName + ".web.secret")
	webExpireTime := viper.GetDuration(moduleName + ".web.expire_time")

	return &HttpConfig{
		HttpAddress:       httpAddress,
		HttpPort:          httpPort,
		ShutdownDelayTime: shutdownDelayTime,
		Web: &webConfig{
			Enable:     webEnable,
			Password:   webPassword,
			Secret:     webSecret,
			ExpireTime: webExpireTime,
		},
	}
}
