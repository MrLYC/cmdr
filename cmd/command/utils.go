package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

func runCommand(fn func(cfg core.Configuration, manager core.CommandManager) error) func(cmd *cobra.Command, args []string) {
	return utils.RunCobraCommandWith(core.CommandProviderDefault, fn)
}
