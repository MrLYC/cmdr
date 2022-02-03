package utils

import (
	"context"
	"fmt"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/spf13/cobra"
)

func ExecuteRunner(factory func(define.Configuration, *CmdrHelper) define.Runner) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := config.Global
		logger := define.Logger
		name := cmd.Name()

		logger.Info("running", map[string]interface{}{
			"command": name,
		})

		ctx := context.WithValue(cmd.Context(), define.ContextKeyConfiguration, cfg)
		helper := NewCmdrHelper(cfg.GetString(config.CfgKeyCmdrRoot))
		runner := factory(cfg, helper)
		ExitWithError(runner.Run(ctx), fmt.Sprintf("command %s failed", name))
	}
}
