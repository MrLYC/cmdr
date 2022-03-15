package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove command from cmdr",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		logger := core.GetLogger()
		name := cfg.GetString(core.CfgKeyXCommandRemoveName)
		version := cfg.GetString(core.CfgKeyXCommandRemoveVersion)

		err := manager.Undefine(name, version)

		if errors.Cause(err) == core.ErrCommandAlreadyActivated {
			logger.Warn("command is already activated, please deactivate it first", map[string]interface{}{
				"name":    name,
				"version": version,
			})
			return nil
		}

		if err != nil {
			return errors.WithMessagef(err, "failed to remove command %s", name)
		}

		logger.Info("command removed", map[string]interface{}{
			"name":    name,
			"version": version,
		})

		return nil
	}),
}

func init() {
	Cmd.AddCommand(removeCmd)
	flags := removeCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXCommandRemoveName, flags.Lookup("name")),
		removeCmd.MarkFlagRequired("name"),

		cfg.BindPFlag(core.CfgKeyXCommandRemoveVersion, flags.Lookup("version")),
		removeCmd.MarkFlagRequired("version"),

		utils.NewDefaultCobraCommandCompleteHelper(removeCmd).RegisterAll(),
	)
}
