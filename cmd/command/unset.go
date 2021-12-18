package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Deactivate a command",
	Run: func(cmd *cobra.Command, args []string) {
		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewNamedCommandsQuerier(simpleCmdFlag.name),
			core.NewCommandDeactivator(),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "deactivate failed")

		define.Logger.Info("unset command", map[string]interface{}{
			"name": simpleCmdFlag.name,
		})
	},
}

func init() {
	Cmd.AddCommand(unsetCmd)

	flags := unsetCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")

	unsetCmd.MarkFlagRequired("name")
}
