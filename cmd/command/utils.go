package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
	"github.com/mrlyc/cmdr/cmdr/utils"
)

func runCommand(fn func(cfg cmdr.Configuration, manager cmdr.CommandManager) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := cmdr.GetConfiguration()

		manager, err := cmdr.NewCommandManager(cmdr.CommandProviderSimple, cfg)
		if err != nil {
			utils.ExitWithError(err, "Failed to create command manager")
		}

		utils.ExitWithError(fn(cfg, manager), "Failed to run command %s", cmd.Name())
	}
}
