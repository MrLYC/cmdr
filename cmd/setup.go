package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

var setupCmdFlag struct {
	doNotInstall      bool
	doNotWriteProfile bool
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		runner := core.NewStepRunner(
			core.NewDirectoryMaker(map[string]string{
				"shims": core.GetShimsDir(),
				"bin":   core.GetBinDir(),
			}),
			core.NewDBClientMaker(),
			core.NewDBMigrator(new(model.Command)),
			core.NewShellProfiler(os.Getenv("SHELL")),
		)

		cmdrLocation, err := os.Executable()
		utils.CheckError(err)

		if !setupCmdFlag.doNotInstall {
			runner.Add(
				core.NewBinaryInstaller(),
				core.NewCommandInstaller(),
			)
		}

		if !setupCmdFlag.doNotWriteProfile {
			runner.Add(
				core.NewShellProfiler(os.Getenv("SHELL")),
			)
		}

		utils.ExitWithError(runner.Run(utils.SetIntoContext(cmd.Context(), map[define.ContextKey]interface{}{
			define.ContextKeyName:           define.Name,
			define.ContextKeyVersion:        define.Version,
			define.ContextKeyLocation:       cmdrLocation,
			define.ContextKeyCommandManaged: true,
		})), "setup failed")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	flags := setupCmd.Flags()
	flags.BoolVar(&setupCmdFlag.doNotInstall, "skip-install", false, "do not install cmdr")
	flags.BoolVar(&setupCmdFlag.doNotWriteProfile, "skip-profile", false, "do not write profile")
}
