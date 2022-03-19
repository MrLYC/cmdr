package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade cmdr",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()
		cfg.Set(core.CfgKeyXUpgradeArgs, []string{"init", "--upgrade"})
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger := core.GetLogger()
		ctx := cmd.Context()
		cfg := core.GetConfiguration()
		releaseName := cfg.GetString(core.CfgKeyXUpgradeRelease)
		assetName := cfg.GetString(core.CfgKeyXUpgradeAsset)
		upgradeArgs := append(cfg.GetStringSlice(core.CfgKeyXUpgradeArgs), args...)

		searcher, err := core.NewCmdrSearcher(core.CmdrSearcherProviderDefault, cfg)
		utils.ExitOnError("getting cmdr searcher", err)

		logger.Info("searching for release", map[string]interface{}{
			"release": releaseName,
		})
		info, err := searcher.GetLatestAsset(ctx, releaseName, assetName)
		utils.ExitOnError("get latest asset url failed", err)

		if info.Version == core.Version {
			logger.Info("cmdr already latest version", map[string]interface{}{
				"version": core.Version,
			})
			return
		}

		utils.ExitOnError("upgrade cmdr failed", utils.UpgradeCmdr(
			ctx, cfg, info.Url, info.Version, upgradeArgs,
		))
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	cfg := core.GetConfiguration()

	flags := upgradeCmd.Flags()
	flags.StringP("release", "r", "latest", "cmdr release tag name")
	flags.StringP("asset", "a", core.Asset, "cmdr release assert name")

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXUpgradeRelease, flags.Lookup("release")),
		cfg.BindPFlag(core.CfgKeyXUpgradeAsset, flags.Lookup("asset")),
	)
}
