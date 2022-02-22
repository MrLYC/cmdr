package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

var configCmdFlag struct {
	list bool
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cmdr configurations",
	Run: func(cmd *cobra.Command, args []string) {
		logger := core.Logger
		cfg := core.GetConfiguration()

		if configCmdFlag.list {
			logger.Debug("listing configurations", map[string]interface{}{
				"path": cfg.ConfigFileUsed(),
			})

			settings := make(map[string]interface{})
			for key, value := range cfg.AllSettings() {
				if !strings.HasPrefix(key, "_") {
					settings[key] = value
				}
			}

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
