package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		return manager.Activate(
			cfg.GetString(core.CfgKeyXCommandUseName),
			cfg.GetString(core.CfgKeyXCommandUseVersion),
		)
	}),
}

func init() {
	Cmd.AddCommand(useCmd)
	flags := useCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",

		cfg.BindPFlag(core.CfgKeyXCommandUseName, flags.Lookup("name")),
		useCmd.MarkFlagRequired("name"),

		cfg.BindPFlag(core.CfgKeyXCommandUseVersion, flags.Lookup("version")),
		useCmd.MarkFlagRequired("version"),
	)
}
