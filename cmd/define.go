package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var defineCmdFlag struct {
	name     string
	version  string
	location string
}

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define",
	Short: "Define command into cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger

		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)

		logger.Debug("defining command", map[string]interface{}{
			"name":     defineCmdFlag.name,
			"version":  defineCmdFlag.version,
			"location": defineCmdFlag.location,
		})

		utils.ExitWithError(
			helper.Define(cmd.Context(), defineCmdFlag.name, defineCmdFlag.version, defineCmdFlag.location),
			"define command failed",
		)

		logger.Info("command defined", map[string]interface{}{
			"name":    defineCmdFlag.name,
			"version": defineCmdFlag.version,
		})

	},
}

func init() {
	rootCmd.AddCommand(defineCmd)

	flags := defineCmd.Flags()
	flags.StringVarP(&defineCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&defineCmdFlag.version, "version", "v", "", "command version")
	flags.StringVarP(&defineCmdFlag.location, "location", "l", "", "command location")

	defineCmd.MarkFlagRequired("name")
	defineCmd.MarkFlagRequired("version")
	defineCmd.MarkFlagRequired("location")
}
