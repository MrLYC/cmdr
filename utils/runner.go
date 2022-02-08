package utils

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
)

func ExecuteRunner(factory func(define.Configuration, *core.Cmdr) define.Runner) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := config.Global
		logger := define.Logger
		name := cmd.Name()

		logger.Info("running", map[string]interface{}{
			"command": name,
		})

		ctx := context.WithValue(cmd.Context(), define.ContextKeyConfiguration, cfg)
		cmdr, err := core.NewCmdr(cfg.GetString(config.CfgKeyCmdrRoot))
		if err != nil {
			ExitWithError(err, "create cmdr failed")
		}
		runner := factory(cfg, cmdr)
		ExitWithError(runner.Run(ctx), fmt.Sprintf("command %s failed", name))
	}
}
