package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/runner"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run:   executeRunner(runner.NewInstallRunner),
}

func init() {
	Cmd.AddCommand(installCmd)
	flags := installCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")
	flags.BoolP("activate", "a", false, "activate command")

	cfg := config.Global
	cfg.BindPFlag(config.CfgKeyCommandInstallName, flags.Lookup("name"))
	cfg.BindPFlag(config.CfgKeyCommandInstallVersion, flags.Lookup("version"))
	cfg.BindPFlag(config.CfgKeyCommandInstallLocation, flags.Lookup("location"))
	cfg.BindPFlag(config.CfgKeyCommandInstallActivate, flags.Lookup("activate"))

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}
