package command

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Deactivate a command",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		logger := core.GetLogger()
		name := cfg.GetString(core.CfgKeyXCommandUnsetName)
		if name == core.Name {
			logger.Error("it is not allowed to unset cmdr itself")
			return nil
		}

		err := manager.Deactivate(name)
		if err != nil {
			return errors.WithMessagef(err, "failed to deactivate command %s", name)
		}

		logger.Info("command deactivated", map[string]interface{}{
			"name": name,
		})

		return nil
	}),
}

func init() {
	Cmd.AddCommand(unsetCmd)
	flags := unsetCmd.Flags()
	flags.StringP("name", "n", "", "command name")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXCommandUnsetName, flags.Lookup("name")),
		unsetCmd.MarkFlagRequired("name"),

		utils.NewDefaultCobraCommandCompleteHelper(unsetCmd).RegisterAll(),
	)
}
