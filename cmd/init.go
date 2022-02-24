package cmd

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

var initCmdFlag struct {
	skipInstall bool
	skipProfile bool
	upgrade     bool
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		logger := core.Logger
		cfg := core.GetConfiguration()
		var errs error

		for _, step := range []string{
			"profile-dir-backup",
			"binary",
			"database-migrator",
			"profile-dir-export",
			"profile-dir-render",
			"profile-injector",
			"cmdr-updater",
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
	rootCmd.AddCommand(initCmd)

	flags := initCmd.Flags()
	flags.BoolVar(&initCmdFlag.upgrade, "upgrade", false, "for upgrade init")
}
