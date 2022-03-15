package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		logger := core.GetLogger()
		name := cfg.GetString(core.CfgKeyXCommandUseName)
		version := cfg.GetString(core.CfgKeyXCommandUseVersion)

		err := manager.Activate(name, version)
		if err != nil {
			return errors.WithMessagef(err, "failed to activate command %s", name)
		}

		logger.Info("command activated", map[string]interface{}{
			"name":    name,
			"version": version,
		})

		return nil
	}),
}

func init() {
	Cmd.AddCommand(useCmd)
	flags := useCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",

		cfg.BindPFlag(core.CfgKeyXCommandUseName, flags.Lookup("name")),
		useCmd.MarkFlagRequired("name"),

		cfg.BindPFlag(core.CfgKeyXCommandUseVersion, flags.Lookup("version")),
		useCmd.MarkFlagRequired("version"),

		utils.NewDefaultCobraCommandCompleteHelper(useCmd).RegisterAll(),
	)
}
