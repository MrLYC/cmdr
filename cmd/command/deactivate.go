package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate a command",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger

		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)

		logger.Debug("deactivating command", map[string]interface{}{
			"name": simpleCmdFlag.name,
		})

		utils.ExitWithError(
			helper.Deactivate(cmd.Context(), simpleCmdFlag.name),
			"deactivate command failed",
		)

		logger.Info("command deactivated", map[string]interface{}{
			"name": simpleCmdFlag.name,
		})
	},
}

func init() {
	Cmd.AddCommand(deactivateCmd)

	flags := deactivateCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")

	deactivateCmd.MarkFlagRequired("name")
}
