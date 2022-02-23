package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// listCmd represents the config command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configurations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()

		logger := core.Logger
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
		utils.ExitOnError("Marshaling settings", err)

		fmt.Printf("%s\n", content)
	},
}

func init() {
	Cmd.AddCommand(listCmd)
}
