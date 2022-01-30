package command

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

var listCmdFlag struct {
	activated bool
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Run: func(cmd *cobra.Command, args []string) {
		runner := operator.NewOperatorRunner(
			operator.NewDBClientMaker(),
			operator.NewFullCommandsQuerier(
				simpleCmdFlag.name, simpleCmdFlag.version, simpleCmdFlag.location, listCmdFlag.activated,
			),
			operator.NewCommandSorter(),
			operator.NewCommandPrinter(os.Stdout),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "list failed")
	},
}

func init() {
	Cmd.AddCommand(listCmd)
	cmdFlagsHelper.declareFlagName(listCmd)
	cmdFlagsHelper.declareFlagVersion(listCmd)
	cmdFlagsHelper.declareFlagLocation(listCmd)
}
