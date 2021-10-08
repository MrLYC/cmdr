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
		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)
		for n, p := range map[string]string{
			"shims": helper.ShimsDir,
			"bin":   helper.BinDir,
		} {
			logger.Info("createing dir", map[string]interface{}{
				"name": n,
				"dir":  p,
			})
			utils.CheckError(define.FS.MkdirAll(p, 0755))
		}

		ctx := cmd.Context()
		logger.Info("creating cmdr database")
		utils.CheckError(client.Schema.Create(ctx))

		if !initCmdFlag.install {
			return
		}

		cmdrPath, err := os.Executable()
		utils.CheckError(err)

		logger.Info("installing cmdr")

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
