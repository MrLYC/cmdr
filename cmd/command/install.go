package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()
		cfg.Set(core.CfgKeyCmdrLinkMode, "default")
	},
	Run: utils.RunCobraCommandWith(core.CommandProviderDownload, func(cfg core.Configuration, manager core.CommandManager) error {
		logger := core.GetLogger()
		name := cfg.GetString(core.CfgKeyXCommandInstallName)
		version := cfg.GetString(core.CfgKeyXCommandInstallVersion)
		location := cfg.GetString(core.CfgKeyXCommandInstallLocation)
		activate := cfg.GetBool(core.CfgKeyXCommandInstallActivate)
		_, err := utils.DefineCmdrCommand(manager, name, version, location, activate)
		if err != nil {
			return errors.WithMessagef(err, "failed to install command %s:%s", name, version)
		}

		logger.Info("command installed", map[string]interface{}{
			"name":    name,
			"version": version,
		})

		return nil
	}),
}

func init() {
	Cmd.AddCommand(InstallCmd)
	flags := InstallCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")
	flags.BoolP("activate", "a", false, "activate command")

	helper := utils.NewDefaultCobraCommandCompleteHelper(InstallCmd)
	cfg := core.GetConfiguration()
	utils.PanicOnError("binding flags",

		cfg.BindPFlag(core.CfgKeyXCommandInstallName, flags.Lookup("name")),
		InstallCmd.MarkFlagRequired("name"),

		cfg.BindPFlag(core.CfgKeyXCommandInstallVersion, flags.Lookup("version")),
		InstallCmd.MarkFlagRequired("version"),

		cfg.BindPFlag(core.CfgKeyXCommandInstallLocation, flags.Lookup("location")),
		InstallCmd.MarkFlagRequired("location"),

		cfg.BindPFlag(core.CfgKeyXCommandInstallActivate, flags.Lookup("activate")),

		helper.RegisterNameFunc(),
		helper.RegisterVersionFunc(),
	)
}
