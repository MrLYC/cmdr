package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/operator"
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
		shimsDir := operator.GetShimsDir()
		binDir := operator.GetBinDir()
		runner := operator.NewOperatorRunner(
			operator.NewDirectoryMaker(map[string]string{
				"shims": shimsDir,
				"bin":   operator.GetBinDir(),
			}),
			operator.NewDBClientMaker(),
			operator.NewDBMigrator(new(model.Command)),
		)

		cmdrLocation, err := os.Executable()
		utils.CheckError(err)

		if !setupCmdFlag.skipInstall && !setupCmdFlag.upgrade {
			runner.Add(
				operator.NewCommandDefiner(shimsDir, define.Name, define.Version, cmdrLocation, true),
				operator.NewBinariesInstaller(shimsDir),
			)
		}

		if !setupCmdFlag.skipProfile {
			runner.Add(
				operator.NewOperatorLoggerWithFields("writing profile"),
				operator.NewShellProfiler(binDir),
			)
		}

		utils.ExitWithError(runner.Run(cmd.Context()), "setup failed")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	flags := setupCmd.Flags()
	flags.BoolVar(&setupCmdFlag.skipInstall, "skip-install", false, "do not install cmdr")
	flags.BoolVar(&setupCmdFlag.skipProfile, "skip-profile", false, "do not write profile")
	flags.BoolVar(&setupCmdFlag.upgrade, "upgrade", false, "for upgrade setup")
}
