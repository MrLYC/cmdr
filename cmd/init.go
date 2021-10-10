package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var initCmdFlag struct {
	install bool
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
		logger.Info("creating cmdr database")

		client := core.GetClient()
		defer utils.CallClose(client)

		utils.CheckError(client.Schema.Create(ctx))

		if !initCmdFlag.install {
			return
		}

		cmdrPath, err := os.Executable()
		utils.CheckError(err)

		logger.Info("installing cmdr")

		helper := core.NewCommandHelper(client)
		command, err := helper.GetCommandByNameAndVersion(ctx, define.Name, define.Version)
		utils.CheckError(err)
		if command == nil {
			utils.CheckError(helper.Install(ctx, define.Name, define.Version, cmdrPath))
			utils.CheckError(helper.Activate(ctx, define.Name, define.Version))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	flags := initCmd.Flags()
	flags.BoolVarP(&initCmdFlag.install, "install", "i", true, "install cmdr")
}
