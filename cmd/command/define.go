package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define",
	Short: "Define command into cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		shimsDir := core.GetShimsDir()

		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewCommandDefiner(shimsDir, simpleCmdFlag.name, simpleCmdFlag.version, simpleCmdFlag.location, false),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "install failed")

		define.Logger.Info("defined command", map[string]interface{}{
			"name":     simpleCmdFlag.name,
			"version":  simpleCmdFlag.version,
			"location": simpleCmdFlag.location,
		})
	},
}

func init() {
	Cmd.AddCommand(defineCmd)
	cmdFlagsHelper.declareFlagName(defineCmd)
	cmdFlagsHelper.declareFlagVersion(defineCmd)

	defineCmd.MarkFlagRequired("name")
	defineCmd.MarkFlagRequired("version")
	defineCmd.MarkFlagRequired("location")
}
