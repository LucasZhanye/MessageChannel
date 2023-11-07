package config

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

// LoadConfig
func LoadConfig(configFile string) {
	// 指定配置文件
	viper.SetConfigFile(configFile)

	// 设置环境变量前缀
	viper.SetEnvPrefix("MESSAGE_CHANNEL")
	// 设置分隔符，应用可以使用.去分隔变量，而保存在环境变量中的是_
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	// 查找并读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config file from %s error: %v", configFile, err)
	}
}

// BindFlag
func BindFlag(cmd *cobra.Command, flags []string) {
	for _, flag := range flags {
		viper.BindPFlag(flag, cmd.Flags().Lookup(flag))
	}
}
