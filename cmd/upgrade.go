package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var upgradeCmdFlag struct {
	release string
	asset   string
	keep    bool
}

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewReleaseSearcher(upgradeCmdFlag.release, upgradeCmdFlag.asset),
			core.NewDownloader(),
			core.NewBinaryInstaller(),
			core.NewCommandInstaller(),
			core.NewBinaryActivator(),
			core.NewCommandDeactivator(),
			core.NewCommandActivator(),
		)

		if !upgradeCmdFlag.keep {
			runner.Add(
				core.NewContextValueSetter(map[define.ContextKey]interface{}{
					define.ContextKeyVersion: define.Version,
				}),
				core.NewCommandListQuerierByNameAndVersion(
					define.Name, define.Version,
				),
				core.NewBinaryRemover(),
				core.NewCommandRemover(),
			)
		}

		utils.ExitWithError(runner.Run(utils.SetIntoContext(cmd.Context(), map[define.ContextKey]interface{}{
			define.ContextKeyName:           define.Name,
			define.ContextKeyCommandManaged: true,
		})), "upgrade failed")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	flags := upgradeCmd.Flags()
	flags.StringVarP(&upgradeCmdFlag.release, "release", "r", "latest", "cmdr release tag name")
	flags.StringVarP(&upgradeCmdFlag.asset, "asset", "a", define.Asset, "cmdr release assert name")
	flags.BoolVarP(&upgradeCmdFlag.keep, "keep", "k", false, "keep the last cmdr version")
}
