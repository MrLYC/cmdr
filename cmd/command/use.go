package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run: func(cmd *cobra.Command, args []string) {
		binDir := core.GetBinDir()
		shimsDir := core.GetShimsDir()

		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewSimpleCommandsQuerier(simpleCmdFlag.name, simpleCmdFlag.version).StrictMode(),
			core.NewCommandDeactivator(),
			core.NewBinariesActivator(binDir, shimsDir),
			core.NewCommandActivator(),
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

	flags := useCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&simpleCmdFlag.version, "version", "v", "", "command version")

	useCmd.MarkFlagRequired("name")
	useCmd.MarkFlagRequired("version")
}
