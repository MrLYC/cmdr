package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/runner"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall command from cmdr",
	Run:   executeRunner(runner.NewUninstallRunner),
}

func init() {
	Cmd.AddCommand(uninstallCmd)
	flags := uninstallCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := config.Global
	cfg.BindPFlag(runner.CfgKeyCommandUninstallName, flags.Lookup("name"))
	cfg.BindPFlag(runner.CfgKeyCommandUninstallVersion, flags.Lookup("version"))

	uninstallCmd.MarkFlagRequired("name")
	uninstallCmd.MarkFlagRequired("version")
}
