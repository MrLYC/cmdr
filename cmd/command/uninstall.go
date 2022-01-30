package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall command from cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		runner := operator.NewOperatorRunner(
			operator.NewDBClientMaker(),
			operator.NewSimpleCommandsQuerier(simpleCmdFlag.name, simpleCmdFlag.version).StrictMode(),
			operator.NewCommandUndefiner(),
			operator.NewBinariesUninstaller(),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "list failed")

		define.Logger.Info("uninstalled command", map[string]interface{}{
			"name":    simpleCmdFlag.name,
			"version": simpleCmdFlag.version,
		})
	},
}

func init() {
	Cmd.AddCommand(uninstallCmd)
	cmdFlagsHelper.declareFlagName(uninstallCmd)
	cmdFlagsHelper.declareFlagVersion(uninstallCmd)

	uninstallCmd.MarkFlagRequired("name")
	uninstallCmd.MarkFlagRequired("version")
}
