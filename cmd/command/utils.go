package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

func runCommand(fn func(cfg core.Configuration, manager core.CommandManager) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()

		manager, err := core.NewCommandManager(core.CommandProviderSimple, cfg)
		if err != nil {
			utils.ExitOnError("Failed to create command manager", err)
		}

		utils.ExitOnError(fmt.Sprintf("Failed to run command %s", cmd.Name()), fn(cfg, manager))
	}
}

func defineCommand(manager core.CommandManager, name string, version string, location string, activate bool) error {
	err := manager.Define(name, version, location)
	if err != nil {
		return err
	}

	if activate {
		return manager.Activate(name, version)
	}

	return nil
}
