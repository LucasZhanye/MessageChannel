package command

import (
	"fmt"
	"messagechannel/internal/pkg/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show MessageChannel version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("current version: %s\n", version.Get())
	},
}
