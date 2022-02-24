package cmd

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

var setupCmdFlag struct {
	skipInstall bool
	skipProfile bool
	upgrade     bool
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		logger := core.Logger
		cfg := core.GetConfiguration()
		var errs error
		for _, step := range []string{
			"profile-dir-backup",
			"binary",
			"database",
			"profile-dir-export",
			"profile-dir-render",
			"profile-injector",
		} {
			logger.Debug("initializing", map[string]interface{}{
				"step": step,
			})
			handler, err := core.NewInitializer(step, cfg)
			if err != nil {
				errs = multierror.Append(errs, errors.WithMessagef(err, "failed to create %s", step))
				continue
			}

			err = handler.Init()
			if err != nil {
				errs = multierror.Append(errs, errors.WithMessagef(err, "failed to initialize %s", step))
				continue
			}
		}

		utils.ExitOnError("Failed to init cmdr", errs)
		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	flags := setupCmd.Flags()
	flags.BoolVar(&setupCmdFlag.skipInstall, "skip-install", false, "do not install cmdr")
	flags.BoolVar(&setupCmdFlag.skipProfile, "skip-profile", false, "do not write profile")
	flags.BoolVar(&setupCmdFlag.upgrade, "upgrade", false, "for upgrade setup")
}
