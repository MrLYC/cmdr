package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

func executeRunner(factory func(define.Configuration) runner.Runner) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := config.Global
		logger := define.Logger
		name := cmd.Name()

		logger.Info("running", map[string]interface{}{
			"command": name,
		})

		ctx := cmd.Context()
		runner := factory(cfg)
		utils.ExitWithError(runner.Run(context.WithValue(ctx, define.ContextKeyConfiguration, cfg)), fmt.Sprintf("command %s failed", name))
	}
}
