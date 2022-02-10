package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define",
	Short: "Define command into cmdr",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		return manager.Define(
			cfg.GetString(core.CfgKeyCommandDefineName),
			cfg.GetString(core.CfgKeyCommandDefineVersion),
			cfg.GetString(core.CfgKeyCommandDefineLocation),
		)
	}),
}

func init() {
	Cmd.AddCommand(defineCmd)
	flags := defineCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")

	cfg := core.GetConfiguration()
	cfg.BindPFlag(core.CfgKeyCommandDefineName, flags.Lookup("name"))
	cfg.BindPFlag(core.CfgKeyCommandDefineVersion, flags.Lookup("version"))
	cfg.BindPFlag(core.CfgKeyCommandDefineLocation, flags.Lookup("location"))

	defineCmd.MarkFlagRequired("name")
	defineCmd.MarkFlagRequired("version")
	defineCmd.MarkFlagRequired("location")
}
