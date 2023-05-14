package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

func runCommand(fn func(cfg core.Configuration, manager core.CommandManager) error) func(cmd *cobra.Command, args []string) {
	return utils.RunCobraCommandWith(core.CommandProviderDefault, fn)
}

func queryCommands(manager core.CommandManager, activate bool, name, version, location string) ([]core.Command, error) {
	query, err := manager.Query()
	if err != nil {
		return nil, err
	}

	if activate {
		query.WithActivated(activate)
	}

	if name != "" {
		query.WithName(name)
	}

	if version != "" {
		query.WithVersion(version)
	}

	if location != "" {
		query.WithLocation(location)
	}

	commands, err := query.All()
	if err != nil {
		return nil, err
	}

	utils.SortCommands(commands)
	return commands, nil
}
