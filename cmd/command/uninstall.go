package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall command from cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewSimpleCommandsQuerier(
				simpleCmdFlag.name, simpleCmdFlag.version,
			),
			core.NewCommandUndefiner(),
			core.NewBinariesUninstaller(),
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

	flags := uninstallCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&simpleCmdFlag.version, "version", "v", "", "command version")

	uninstallCmd.MarkFlagRequired("name")
	uninstallCmd.MarkFlagRequired("version")
}
