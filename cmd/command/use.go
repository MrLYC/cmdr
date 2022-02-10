package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		return manager.Activate(
			cfg.GetString(core.CfgKeyCommandUseName),
			cfg.GetString(core.CfgKeyCommandUseVersion),
		)
	}),
}

func init() {
	Cmd.AddCommand(useCmd)
	flags := useCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := core.GetConfiguration()
	cfg.BindPFlag(core.CfgKeyCommandUseName, flags.Lookup("name"))
	cfg.BindPFlag(core.CfgKeyCommandUseVersion, flags.Lookup("version"))

	useCmd.MarkFlagRequired("name")
	useCmd.MarkFlagRequired("version")
}
