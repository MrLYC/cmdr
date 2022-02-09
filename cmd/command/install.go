package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run: runCommand(func(cfg cmdr.Configuration, manager cmdr.CommandManager) error {
		name := cfg.GetString(cmdr.CfgKeyCommandInstallName)
		version := cfg.GetString(cmdr.CfgKeyCommandInstallVersion)

		err := manager.Define(
			name, version,
			cfg.GetString(cmdr.CfgKeyCommandInstallLocation),
		)
		if err != nil {
			return err
		}

		if cfg.GetBool(cmdr.CfgKeyCommandInstallActivate) {
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

	cfg := cmdr.GetConfiguration()
	cfg.BindPFlag(cmdr.CfgKeyCommandInstallName, flags.Lookup("name"))
	cfg.BindPFlag(cmdr.CfgKeyCommandInstallVersion, flags.Lookup("version"))
	cfg.BindPFlag(cmdr.CfgKeyCommandInstallLocation, flags.Lookup("location"))
	cfg.BindPFlag(cmdr.CfgKeyCommandInstallActivate, flags.Lookup("activate"))

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}
