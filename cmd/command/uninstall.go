package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall command from cmdr",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		return manager.Undefine(
			cfg.GetString(core.CfgKeyXCommandUninstallName),
			cfg.GetString(core.CfgKeyXCommandUninstallVersion),
		)
	}),
}

func init() {
	Cmd.AddCommand(uninstallCmd)
	flags := uninstallCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXCommandUninstallName, flags.Lookup("name")),
		cfg.BindPFlag(core.CfgKeyXCommandUninstallVersion, flags.Lookup("version")),
		uninstallCmd.MarkFlagRequired("name"),
		uninstallCmd.MarkFlagRequired("version"),
	)
}
