package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define",
	Short: "Define command into cmdr",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		return manager.Define(
			cfg.GetString(core.CfgKeyXCommandDefineName),
			cfg.GetString(core.CfgKeyXCommandDefineVersion),
			cfg.GetString(core.CfgKeyXCommandDefineLocation),
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
	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXCommandDefineName, flags.Lookup("name")),
		cfg.BindPFlag(core.CfgKeyXCommandDefineVersion, flags.Lookup("version")),
		cfg.BindPFlag(core.CfgKeyXCommandDefineLocation, flags.Lookup("location")),
		defineCmd.MarkFlagRequired("name"),
		defineCmd.MarkFlagRequired("version"),
		defineCmd.MarkFlagRequired("location"),
	)
}
