package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall command from cmdr",
	Run: runCommand(func(cfg cmdr.Configuration, manager cmdr.CommandManager) error {
		return manager.Undefine(
			cfg.GetString(cmdr.CfgKeyCommandUninstallName),
			cfg.GetString(cmdr.CfgKeyCommandUninstallVersion),
		)
	}),
}

func init() {
	Cmd.AddCommand(uninstallCmd)
	flags := uninstallCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := cmdr.GetConfiguration()
	cfg.BindPFlag(cmdr.CfgKeyCommandUninstallName, flags.Lookup("name"))
	cfg.BindPFlag(cmdr.CfgKeyCommandUninstallVersion, flags.Lookup("version"))

	uninstallCmd.MarkFlagRequired("name")
	uninstallCmd.MarkFlagRequired("version")
}
