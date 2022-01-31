package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

func executeRunner(factory func(define.Configuration) runner.Runner) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := define.Config
		logger := define.Logger
		name := cmd.Name()

		logger.Info("running", map[string]interface{}{
			"command": name,
		})

		runner := factory(cfg)
		utils.ExitWithError(runner.Run(cmd.Context()), fmt.Sprintf("command %s failed", name))
	}
}
