package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		name := cfg.GetString(core.CfgKeyXCommandInstallName)
		version := cfg.GetString(core.CfgKeyXCommandInstallVersion)

		err := manager.Define(
			name, version,
			cfg.GetString(core.CfgKeyXCommandInstallLocation),
		)
		if err != nil {
			return err
		}

		if cfg.GetBool(core.CfgKeyXCommandInstallActivate) {
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
	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXCommandInstallName, flags.Lookup("name")),
		cfg.BindPFlag(core.CfgKeyXCommandInstallVersion, flags.Lookup("version")),
		cfg.BindPFlag(core.CfgKeyXCommandInstallLocation, flags.Lookup("location")),
		cfg.BindPFlag(core.CfgKeyXCommandInstallActivate, flags.Lookup("activate")),
		installCmd.MarkFlagRequired("name"),
		installCmd.MarkFlagRequired("version"),
		installCmd.MarkFlagRequired("location"),
	)
}
