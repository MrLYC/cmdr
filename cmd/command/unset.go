package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Deactivate a command",
	Run: func(cmd *cobra.Command, args []string) {
		runner := operator.NewOperatorRunner(
			operator.NewDBClientMaker(),
			operator.NewNamedCommandsQuerier(simpleCmdFlag.name),
			operator.NewCommandDeactivator(),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "deactivate failed")

		define.Logger.Info("unset command", map[string]interface{}{
			"name": simpleCmdFlag.name,
		})
	},
}

func init() {
	Cmd.AddCommand(unsetCmd)
	cmdFlagsHelper.declareFlagName(unsetCmd)

	unsetCmd.MarkFlagRequired("name")
}
