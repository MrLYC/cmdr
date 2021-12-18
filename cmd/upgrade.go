package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var upgradeCmdFlag struct {
	release   string
	asset     string
	keep      bool
	skipSetup bool
}

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger
		runner := core.NewStepRunner()
		shimsDir := core.GetShimsDir()
		binDir := core.GetBinDir()
		cmdrLocation, err := os.Executable()
		utils.CheckError(err)

		if !upgradeCmdFlag.skipSetup {
			runArgs := []string{"setup", "--upgrade"}
			runArgs = append(runArgs, args...)
			runner.Add(core.NewUpgradeSetupRunner(runArgs...))
		}

		runner.Add(
			core.NewDBClientMaker(),
			core.NewCommandDefiner(shimsDir, define.Name, define.Version, cmdrLocation, true),
			core.NewReleaseSearcher(upgradeCmdFlag.release, upgradeCmdFlag.asset),
			core.NewDownloader(),
			core.NewBinariesInstaller(shimsDir),
		)

		if !upgradeCmdFlag.keep {
			runner.Add(
				core.NewCommandDeactivator(),
				core.NewBinariesActivator(binDir),
				core.NewCommandActivator(),
				core.NewSimpleCommandsQuerier(
					define.Name, define.Version,
				),
				core.NewCommandUndefiner(),
				core.NewBinariesUninstaller(),
			)
		}

		utils.ExitWithError(runner.Run(cmd.Context()), "upgrade failed")

		logger.Info("upgraded command", map[string]interface{}{
			"name": define.Name,
		})
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	flags := upgradeCmd.Flags()
	flags.StringVarP(&upgradeCmdFlag.release, "release", "r", "latest", "cmdr release tag name")
	flags.StringVarP(&upgradeCmdFlag.asset, "asset", "a", define.Asset, "cmdr release assert name")
	flags.BoolVarP(&upgradeCmdFlag.keep, "keep", "k", false, "keep the last cmdr version")
	flags.BoolVar(&upgradeCmdFlag.skipSetup, "skip-setup", false, "do not setup after cmdr installed")
}
