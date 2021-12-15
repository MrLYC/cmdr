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
	skipInstall bool
	skipProfile bool
	upgrade     bool
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		shimsDir := core.GetShimsDir()
		runner := core.NewStepRunner(
			core.NewDirectoryMaker(map[string]string{
				"shims": shimsDir,
				"bin":   core.GetBinDir(),
			}),
			core.NewDBClientMaker(),
			core.NewDBMigrator(new(model.Command)),
			core.NewShellProfiler(os.Getenv("SHELL")),
		)

		cmdrLocation, err := os.Executable()
		utils.CheckError(err)

		if !setupCmdFlag.skipInstall && !setupCmdFlag.upgrade {
			runner.Add(
				core.NewStepLoggerWithFields("installing command", define.ContextKeyName, define.ContextKeyVersion),
				core.NewBinariesInstaller(shimsDir),
				core.NewCommandDefiner(shimsDir),
			)
		}

		if !setupCmdFlag.skipProfile {
			runner.Add(
				core.NewStepLoggerWithFields("writing profile"),
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
	flags.BoolVar(&setupCmdFlag.skipInstall, "skip-install", false, "do not install cmdr")
	flags.BoolVar(&setupCmdFlag.skipProfile, "skip-profile", false, "do not write profile")
	flags.BoolVar(&setupCmdFlag.upgrade, "upgrade", false, "for upgrade setup")
}
