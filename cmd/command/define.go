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
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()
		cfg.Set(core.CfgKeyCmdrLinkMode, "link")
	},
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		return defineCommand(
			manager,
			cfg.GetString(core.CfgKeyXCommandDefineName),
			cfg.GetString(core.CfgKeyXCommandDefineVersion),
			cfg.GetString(core.CfgKeyXCommandDefineLocation),
			cfg.GetBool(core.CfgKeyXCommandDefineActivate),
		)
	}),
}

func init() {
	Cmd.AddCommand(defineCmd)
	flags := defineCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")
	flags.BoolP("activate", "a", false, "activate command")

	cfg := core.GetConfiguration()
	utils.PanicOnError("binding flags",

		cfg.BindPFlag(core.CfgKeyXCommandDefineName, flags.Lookup("name")),
		defineCmd.MarkFlagRequired("name"),

		cfg.BindPFlag(core.CfgKeyXCommandDefineVersion, flags.Lookup("version")),
		defineCmd.MarkFlagRequired("version"),

		cfg.BindPFlag(core.CfgKeyXCommandDefineLocation, flags.Lookup("location")),
		defineCmd.MarkFlagRequired("location"),

		cfg.BindPFlag(core.CfgKeyXCommandDefineActivate, flags.Lookup("activate")),
	)
}
