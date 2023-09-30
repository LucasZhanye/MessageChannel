package command

import (
	"messagechannel/internal/pkg/core"
	"messagechannel/internal/pkg/core/server"

	"messagechannel/pkg/logger"

	"github.com/spf13/cobra"
)

var configFile string // config file path such as config/config.yaml
var flagSet []string  // all flag's name

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "MessageChannel root cmd",
	Run: func(cmd *cobra.Command, args []string) {
		// binding flag
		core.BindFlag(cmd, flagSet)
		// load config
		core.LoadConfig(configFile)

		// init Logger
		log, _ := logger.InitDefaultLog(logger.Config{Level: "info"})

		// new Node
		n := core.NewNode(log)

		go n.HandleSignals()

		// new Server
		s := server.New(n)

		// start Server
		s.Run()
	},
}

// InitCmd init command, add flag and child command here
func InitCmd() *cobra.Command {

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "config file path")
	rootCmd.Flags().StringP("test", "t", "abc", "test string")

	flagSet = append(flagSet, "test")
	rootCmd.AddCommand(versionCmd)

	return rootCmd
}
