package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define",
	Short: "Define command into cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		runner := runner.NewDefineRunner(define.Config)
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

	flags := installCmd.Flags()
	cfg := define.Config
	cfg.BindPFlag(runner.CfgKeyCommandDefineName, flags.Lookup("name"))
	cfg.BindPFlag(runner.CfgKeyCommandDefineVersion, flags.Lookup("version"))
	cfg.BindPFlag(runner.CfgKeyCommandDefineLocation, flags.Lookup("location"))

	defineCmd.MarkFlagRequired("name")
	defineCmd.MarkFlagRequired("version")
	defineCmd.MarkFlagRequired("location")
}
