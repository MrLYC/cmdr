package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// DefineCmd represents the define command
var DefineCmd = &cobra.Command{
	Use:   "define",
	Short: "Define command into cmdr",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()
		cfg.Set(core.CfgKeyCmdrLinkMode, "link")
	},
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		logger := core.GetLogger()
		name := cfg.GetString(core.CfgKeyXCommandDefineName)
		version := cfg.GetString(core.CfgKeyXCommandDefineVersion)
		location := cfg.GetString(core.CfgKeyXCommandDefineLocation)
		activate := cfg.GetBool(core.CfgKeyXCommandDefineActivate)
		_, err := utils.DefineCmdrCommand(manager, name, version, location, activate)
		if err != nil {
			return errors.WithMessagef(err, "failed to define command %s:%s", name, version)
		}

		logger.Info("command defined", map[string]interface{}{
			"name":    name,
			"version": version,
		})
		return nil
	}),
}

func init() {
	Cmd.AddCommand(DefineCmd)
	flags := DefineCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")
	flags.BoolP("activate", "a", false, "activate command")

	helper := utils.NewDefaultCobraCommandCompleteHelper(DefineCmd)
	cfg := core.GetConfiguration()
	utils.PanicOnError("binding flags",

		cfg.BindPFlag(core.CfgKeyXCommandDefineName, flags.Lookup("name")),
		DefineCmd.MarkFlagRequired("name"),

		cfg.BindPFlag(core.CfgKeyXCommandDefineVersion, flags.Lookup("version")),
		DefineCmd.MarkFlagRequired("version"),

		cfg.BindPFlag(core.CfgKeyXCommandDefineLocation, flags.Lookup("location")),
		DefineCmd.MarkFlagRequired("location"),

		cfg.BindPFlag(core.CfgKeyXCommandDefineActivate, flags.Lookup("activate")),

		helper.RegisterNameFunc(),
		helper.RegisterVersionFunc(),
	)
}
