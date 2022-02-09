package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run: runCommand(func(cfg cmdr.Configuration, manager cmdr.CommandManager) error {
		return manager.Activate(
			cfg.GetString(cmdr.CfgKeyCommandUseName),
			cfg.GetString(cmdr.CfgKeyCommandUseVersion),
		)
	}),
}

func init() {
	Cmd.AddCommand(useCmd)
	flags := useCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := cmdr.GetConfiguration()
	cfg.BindPFlag(cmdr.CfgKeyCommandUseName, flags.Lookup("name"))
	cfg.BindPFlag(cmdr.CfgKeyCommandUseVersion, flags.Lookup("version"))

	useCmd.MarkFlagRequired("name")
	useCmd.MarkFlagRequired("version")
}
