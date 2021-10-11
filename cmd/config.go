package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var configCmdFlag struct {
	list bool
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cmdr configurations",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger

		if configCmdFlag.list {
			logger.Debug("listing configurations", map[string]interface{}{
				"path": define.Configuration.ConfigFileUsed(),
			})

			settings := define.Configuration.AllSettings()
			content, err := yaml.Marshal(settings)
			utils.ExitWithError(err, "marshaling settings")

			fmt.Printf("%s\n", content)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	flags := configCmd.Flags()
	flags.BoolVarP(&configCmdFlag.list, "list", "l", false, "list config")
}
