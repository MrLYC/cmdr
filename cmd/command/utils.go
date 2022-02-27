package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

func runCommand(fn func(cfg core.Configuration, manager core.CommandManager) error) func(cmd *cobra.Command, args []string) {
	return runCommandWith(core.CommandProviderDefault, fn)
}

func runCommandWith(provider core.CommandProvider, fn func(cfg core.Configuration, manager core.CommandManager) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()

		manager, err := core.NewCommandManager(provider, cfg)
		if err != nil {
			utils.ExitOnError("Failed to create command manager", err)
		}

		defer utils.CallClose(manager)

		utils.ExitOnError(fmt.Sprintf("Failed to run command %s", cmd.Name()), fn(cfg, manager))
	}
}
