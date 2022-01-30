package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run: func(cmd *cobra.Command, args []string) {
		binDir := operator.GetBinDir()
		shimsDir := operator.GetShimsDir()

		runner := operator.NewOperatorRunner(
			operator.NewDBClientMaker(),
			operator.NewSimpleCommandsQuerier(simpleCmdFlag.name, simpleCmdFlag.version).StrictMode(),
			operator.NewCommandDeactivator(),
			operator.NewBinariesActivator(binDir, shimsDir),
			operator.NewCommandActivator(),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "activate failed")

		define.Logger.Info("used command", map[string]interface{}{
			"name":    simpleCmdFlag.name,
			"version": simpleCmdFlag.version,
		})
	},
}

func init() {
	Cmd.AddCommand(useCmd)
	cmdFlagsHelper.declareFlagName(useCmd)
	cmdFlagsHelper.declareFlagVersion(useCmd)

	useCmd.MarkFlagRequired("name")
	useCmd.MarkFlagRequired("version")
}
