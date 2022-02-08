package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/runner"
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
		cfg := config.GetGlobalConfiguration()
		cmdr, err := core.NewCmdr(cfg.GetString(config.CfgKeyCmdrRoot))
		if err != nil {
			utils.ExitWithError(err, "create cmdr failed")
		}

		runner := runner.New(
			operator.NewCommandDefiner(cmdr, define.Name, define.Version, upgradeCmdFlag.location),
		)

		if upgradeCmdFlag.location == "" {
			runner.Add(
				operator.NewReleaseSearcher(upgradeCmdFlag.release, upgradeCmdFlag.asset),
				operator.NewDownloader(),
			)
		}

		runner.Add(
			operator.NewBinariesInstaller(cmdr, true),
		)

		if !upgradeCmdFlag.keep {
			runner.Add(
				operator.NewCommandDeactivator(cmdr),
				operator.NewBinariesActivator(cmdr),
				operator.NewCommandActivator(cmdr),
				operator.NewNamedCommandsQuerier(define.Name),
				operator.NewCommandUndefiner(cmdr),
				operator.NewBinariesUninstaller(cmdr),
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

		binPath, err := cmdr.BinaryManager.BinManager.RealPath(define.Name)
		if err != nil {
			utils.ExitWithError(err, "get binary path failed")
		}
		utils.ExitWithError(utils.WaitProcess(ctx, binPath, runArgs), "setup failed")
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
