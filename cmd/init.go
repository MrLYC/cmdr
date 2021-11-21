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
		runner := core.NewStepRunner(
			core.NewDirectoryMaker(map[string]string{
				"shims": core.GetShimsDir(),
				"bin":   core.GetBinDir(),
			}),
			core.NewDBClientMaker(),
			core.NewDBMigrator(new(model.Command)),
		)

		cmdrLocation, err := os.Executable()
		utils.CheckError(err)

		if !initCmdFlag.doNotInstall {
			runner.Add(
				core.NewBinaryInstaller(),
				core.NewCommandInstaller(),
			)
		}

		utils.ExitWithError(runner.Run(utils.SetIntoContext(cmd.Context(), map[define.ContextKey]interface{}{
			define.ContextKeyName:           define.Name,
			define.ContextKeyVersion:        define.Version,
			define.ContextKeyLocation:       cmdrLocation,
			define.ContextKeyCommandManaged: true,
		})), "init failed")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	flags := initCmd.Flags()
	flags.BoolVar(&initCmdFlag.doNotInstall, "do-not-install-cmdr", false, "do not install cmdr")
}
