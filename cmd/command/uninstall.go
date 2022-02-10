package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall command from cmdr",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		return manager.Undefine(
			cfg.GetString(core.CfgKeyCommandUninstallName),
			cfg.GetString(core.CfgKeyCommandUninstallVersion),
		)
	}),
}

func init() {
	Cmd.AddCommand(uninstallCmd)
	flags := uninstallCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := core.GetConfiguration()
	cfg.BindPFlag(core.CfgKeyCommandUninstallName, flags.Lookup("name"))
	cfg.BindPFlag(core.CfgKeyCommandUninstallVersion, flags.Lookup("version"))

	uninstallCmd.MarkFlagRequired("name")
	uninstallCmd.MarkFlagRequired("version")
}
