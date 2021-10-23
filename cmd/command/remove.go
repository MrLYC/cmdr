package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

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
			"name":    simpleCmdFlag.name,
			"version": simpleCmdFlag.version,
		})
		utils.ExitWithError(
			helper.Remove(cmd.Context(), simpleCmdFlag.name, simpleCmdFlag.version),
			"remove command failed",
		)
		logger.Info("command removed", map[string]interface{}{
			"name":    simpleCmdFlag.name,
			"version": simpleCmdFlag.version,
		})
	},
}

func init() {
	Cmd.AddCommand(removeCmd)

	flags := removeCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&simpleCmdFlag.version, "version", "v", "", "command version")

	removeCmd.MarkFlagRequired("name")
	removeCmd.MarkFlagRequired("version")
}
