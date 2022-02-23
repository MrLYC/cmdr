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
