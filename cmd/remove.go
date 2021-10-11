package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var removeCmdFlag struct {
	name    string
	version string
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove command from cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger
		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)

		logger.Debug("removing command", map[string]interface{}{
			"name":    removeCmdFlag.name,
			"version": removeCmdFlag.version,
		})
		utils.ExitWithError(
			helper.Remove(cmd.Context(), removeCmdFlag.name, removeCmdFlag.version),
			"remove command failed",
		)
		logger.Info("command removed", map[string]interface{}{
			"name":    removeCmdFlag.name,
			"version": removeCmdFlag.version,
		})
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	flags := removeCmd.Flags()
	flags.StringVarP(&removeCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&removeCmdFlag.version, "version", "v", "", "command version")

	removeCmd.MarkFlagRequired("name")
	removeCmd.MarkFlagRequired("version")
}
