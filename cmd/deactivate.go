package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var deactivateCmdFlag struct {
	name string
}

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
			"name": deactivateCmdFlag.name,
		})

		utils.ExitWithError(
			helper.Deactivate(cmd.Context(), deactivateCmdFlag.name),
			"deactivate command failed",
		)

		logger.Info("command deactivated", map[string]interface{}{
			"name": deactivateCmdFlag.name,
		})
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)

	flags := deactivateCmd.Flags()
	flags.StringVarP(&deactivateCmdFlag.name, "name", "n", "", "command name")

	deactivateCmd.MarkFlagRequired("name")
}
