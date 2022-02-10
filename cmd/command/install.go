package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		name := cfg.GetString(core.CfgKeyCommandInstallName)
		version := cfg.GetString(core.CfgKeyCommandInstallVersion)

		err := manager.Define(
			name, version,
			cfg.GetString(core.CfgKeyCommandInstallLocation),
		)
		if err != nil {
			return err
		}

		if cfg.GetBool(core.CfgKeyCommandInstallActivate) {
			return manager.Activate(name, version)
		}

		return nil
	}),
}

func init() {
	Cmd.AddCommand(installCmd)
	flags := installCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")
	flags.BoolP("activate", "a", false, "activate command")

	cfg := core.GetConfiguration()
	cfg.BindPFlag(core.CfgKeyCommandInstallName, flags.Lookup("name"))
	cfg.BindPFlag(core.CfgKeyCommandInstallVersion, flags.Lookup("version"))
	cfg.BindPFlag(core.CfgKeyCommandInstallLocation, flags.Lookup("location"))
	cfg.BindPFlag(core.CfgKeyCommandInstallActivate, flags.Lookup("activate"))

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}
