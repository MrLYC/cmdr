package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/runner"
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
		cfg := config.GetGlobalConfiguration()
		helper := utils.NewCmdrHelper(cfg.GetString(config.CfgKeyCmdrRoot))
		shimsDir := helper.GetShimsDir()
		runner := runner.New(
			operator.NewDirectoryMaker(map[string]string{
				"shims": shimsDir,
				"bin":   helper.GetBinDir(),
			}),
			operator.NewDBClientMaker(helper),
			operator.NewDBMigrator(new(model.Command)),
		)

		cmdrLocation, err := os.Executable()
		utils.CheckError(err)

		if !setupCmdFlag.skipInstall && !setupCmdFlag.upgrade {
			runner.Add(
				operator.NewCommandDefiner(define.Name, define.Version, cmdrLocation, true, helper),
				operator.NewBinariesInstaller(helper),
			)
		}

		if !setupCmdFlag.skipProfile {
			runner.Add(
				operator.NewOperatorLoggerWithFields("writing profile"),
				operator.NewShellProfiler(helper),
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
