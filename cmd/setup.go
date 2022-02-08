package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
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
		cmdr, err := core.NewCmdr(cfg.GetString(config.CfgKeyCmdrRoot))
		if err != nil {
			utils.ExitWithError(err, "create cmdr failed")
		}

		runner := runner.NewMigrateRunner(cfg, cmdr)

		cmdrLocation, err := os.Executable()
		utils.CheckError(err)

		if !setupCmdFlag.skipInstall && !setupCmdFlag.upgrade {
			runner.Add(
				operator.NewCommandDefiner(cmdr, define.Name, define.Version, cmdrLocation),
				operator.NewBinariesInstaller(cmdr, true),
			)
		}

		if !setupCmdFlag.skipProfile {
			runner.Add(
				operator.NewOperatorLoggerWithFields("writing profile"),
				operator.NewShellProfiler(cmdr),
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
