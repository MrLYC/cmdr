package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var upgradeCmdFlag struct {
	release   string
	asset     string
	location  string
	keep      bool
	skipSetup bool
}

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		logger := define.Logger
		shimsDir := core.GetShimsDir()
		binDir := core.GetBinDir()
		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewCommandDefiner(shimsDir, define.Name, define.Version, upgradeCmdFlag.location, true),
		)

		if upgradeCmdFlag.location == "" {
			runner.Add(
				core.NewReleaseSearcher(upgradeCmdFlag.release, upgradeCmdFlag.asset),
				core.NewDownloader(),
			)
		}

		runner.Add(
			core.NewBinariesInstaller(shimsDir),
		)

		if !upgradeCmdFlag.keep {
			runner.Add(
				core.NewCommandDeactivator(),
				core.NewBinariesActivator(binDir, shimsDir),
				core.NewCommandActivator(),
				core.NewNamedCommandsQuerier(define.Name),
				core.NewCommandUndefiner(),
				core.NewBinariesUninstaller(),
			)
		}

		utils.ExitWithError(runner.Run(ctx), "upgrade failed")

		logger.Info("upgraded command", map[string]interface{}{
			"name": define.Name,
		})

		if upgradeCmdFlag.skipSetup {
			return
		}

		runArgs := []string{"setup", "--upgrade"}
		cmd.PersistentFlags().Visit(func(f *pflag.Flag) {
			runArgs = append(runArgs, fmt.Sprintf("--%s=%s", f.Name, f.Value.String()))
		})
		runArgs = append(runArgs, args...)

		logger.Info("setup command", map[string]interface{}{
			"args": runArgs,
		})

		utils.ExitWithError(utils.WaitProcess(ctx, filepath.Join(core.GetBinDir(), define.Name), runArgs), "setup failed")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	flags := upgradeCmd.Flags()
	flags.StringVarP(&upgradeCmdFlag.release, "release", "r", "latest", "cmdr release tag name")
	flags.StringVarP(&upgradeCmdFlag.asset, "asset", "a", define.Asset, "cmdr release assert name")
	flags.StringVarP(&upgradeCmdFlag.location, "location", "l", "", "cmdr binary local location")
	flags.BoolVarP(&upgradeCmdFlag.keep, "keep", "k", false, "keep the last cmdr version")
	flags.BoolVar(&upgradeCmdFlag.skipSetup, "skip-setup", false, "do not setup after cmdr installed")
}
