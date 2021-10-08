package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cmdr configurations",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		isListMode, err := flags.GetBool("list")
		utils.CheckError(err)

		if isListMode {
			settings := define.Configuration.AllSettings()
			content, err := yaml.Marshal(settings)
			utils.CheckError(err)

			fmt.Printf("%s\n", content)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	flags := configCmd.Flags()
	flags.BoolP("list", "l", false, "list config")
}
