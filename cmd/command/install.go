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
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()
		cfg.Set(core.CfgKeyCmdrLinkMode, "default")
	},
	Run: utils.RunCobraCommandWith(core.CommandProviderDownload, func(cfg core.Configuration, manager core.CommandManager) error {
		_, err := utils.DefineCmdrCommand(
			manager,
			cfg.GetString(core.CfgKeyXCommandInstallName),
			cfg.GetString(core.CfgKeyXCommandInstallVersion),
			cfg.GetString(core.CfgKeyXCommandInstallLocation),
			cfg.GetBool(core.CfgKeyXCommandInstallActivate),
		)
		return err
	}),
}

func init() {
	Cmd.AddCommand(installCmd)
	flags := installCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")
	flags.BoolP("activate", "a", false, "activate command")

	helper := utils.NewDefaultCobraCommandCompleteHelper(installCmd)
	cfg := core.GetConfiguration()
	utils.PanicOnError("binding flags",

		cfg.BindPFlag(core.CfgKeyXCommandInstallName, flags.Lookup("name")),
		installCmd.MarkFlagRequired("name"),

		cfg.BindPFlag(core.CfgKeyXCommandInstallVersion, flags.Lookup("version")),
		installCmd.MarkFlagRequired("version"),

		cfg.BindPFlag(core.CfgKeyXCommandInstallLocation, flags.Lookup("location")),
		installCmd.MarkFlagRequired("location"),

		cfg.BindPFlag(core.CfgKeyXCommandInstallActivate, flags.Lookup("activate")),

		helper.RegisterNameFunc(),
		helper.RegisterVersionFunc(),
	)
}
