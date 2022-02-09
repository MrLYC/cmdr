package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define",
	Short: "Define command into cmdr",
	Run: runCommand(func(cfg cmdr.Configuration, manager cmdr.CommandManager) error {
		return manager.Define(
			cfg.GetString(cmdr.CfgKeyCommandDefineName),
			cfg.GetString(cmdr.CfgKeyCommandDefineVersion),
			cfg.GetString(cmdr.CfgKeyCommandDefineLocation),
		)
	}),
}

func init() {
	Cmd.AddCommand(defineCmd)
	flags := defineCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")

	cfg := cmdr.GetConfiguration()
	cfg.BindPFlag(cmdr.CfgKeyCommandDefineName, flags.Lookup("name"))
	cfg.BindPFlag(cmdr.CfgKeyCommandDefineVersion, flags.Lookup("version"))
	cfg.BindPFlag(cmdr.CfgKeyCommandDefineLocation, flags.Lookup("location"))

	defineCmd.MarkFlagRequired("name")
	defineCmd.MarkFlagRequired("version")
	defineCmd.MarkFlagRequired("location")
}
