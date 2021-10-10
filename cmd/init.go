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
	doNotinstall bool
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
			utils.CheckError(define.FS.MkdirAll(p, 0755))
		}

		ctx := cmd.Context()

		client := core.GetClient()
		defer utils.CallClose(client)

		logger.Info("creating cmdr database")
		utils.CheckError(client.Migrate(
			new(model.Command),
		))

		if initCmdFlag.doNotinstall {
			return
		}

		cmdrPath, err := os.Executable()
		utils.CheckError(err)

		logger.Info("installing cmdr")

		helper := core.NewCommandHelper(client)
		installed, err := helper.Upgrade(ctx, define.Version, cmdrPath)
		utils.CheckError(err)

		if installed {
			logger.Info("cmdr already installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	flags := initCmd.Flags()
	flags.BoolVar(&initCmdFlag.doNotinstall, "do-not-install-cmdr", false, "do not install cmdr")
}
