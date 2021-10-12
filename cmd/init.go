package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

var initCmdFlag struct {
	doNotInstall bool
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial cmdr environment",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger
		for n, p := range map[string]string{
			"shims": core.GetShimsDir(),
			"bin":   core.GetBinDir(),
		} {
			logger.Info("createing dir", map[string]interface{}{
				"name": n,
				"dir":  p,
			})
			utils.ExitWithError(
				define.FS.MkdirAll(p, 0755),
				"making dir %s failed", n,
			)
		}

		ctx := cmd.Context()

		client := core.GetClient()
		defer utils.CallClose(client)

		logger.Debug("migrating database")
		utils.ExitWithError(
			client.Migrate(
				new(model.Command),
			),
			"migrate failed",
		)
		logger.Info("database migrated")

		if initCmdFlag.doNotInstall {
			return
		}

		cmdrPath, err := os.Executable()
		utils.CheckError(err)

		helper := core.NewCommandHelper(client)

		logger.Debug("installing cmdr", map[string]interface{}{
			"version": define.Version,
			"path":    cmdrPath,
		})
		installed, err := helper.Upgrade(ctx, define.Version, cmdrPath)
		utils.ExitWithError(err, "cmdr install failed")

		if installed {
			logger.Info("cmdr already installed")
		} else {
			logger.Info("cmdr installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	flags := initCmd.Flags()
	flags.BoolVar(&initCmdFlag.doNotInstall, "do-not-install-cmdr", false, "do not install cmdr")
}
