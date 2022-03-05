package cmd

import (
	"strings"

	"github.com/google/go-github/v39/github"
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
		logger := core.Logger
		ctx := cmd.Context()
		cfg := core.GetConfiguration()
		releaseTag := cfg.GetString(core.CfgKeyXUpgradeRelease)
		assetName := cfg.GetString(core.CfgKeyXUpgradeAsset)
		upgradeArgs := cfg.GetStringSlice(core.CfgKeyXUpgradeArgs)
		githubClient := github.NewClient(nil)

		logger.Info("searching for release", map[string]interface{}{
			"release": releaseTag,
		})
		release, err := utils.GetCmdrRelease(ctx, githubClient.Repositories, releaseTag)
		utils.ExitOnError("get release failed", err)

		if strings.Contains(release.GetTagName(), core.Version) {
			logger.Info("cmdr already latest version", map[string]interface{}{
				"version": core.Version,
			})
			return
		}

		asset, err := utils.SearchReleaseAsset(ctx, assetName, release)
		utils.ExitOnError("search release asset failed", err)

		logger.Info("release asset found", map[string]interface{}{
			"release": releaseTag,
			"asset":   asset.GetName(),
		})
		utils.ExitOnError("upgrade cmdr failed", utils.UpgradeCmdr(
			ctx, cfg, asset.GetBrowserDownloadURL(),
			release.GetTagName(), upgradeArgs,
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
