package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
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
		location := cfg.GetString(core.CfgKeyXUpgradeLocation)
		upgradeArgs := cfg.GetStringSlice(core.CfgKeyXUpgradeArgs)

		downloadFromGithub := func() (string, bool) {
			release, err := utils.GetCMDRRelease(ctx, releaseTag)
			if err != nil {
				utils.ExitOnError("search cmdr release failed", err)
			}

			releaseName := release.GetName()
			if strings.Contains(releaseName, core.Version) {
				logger.Info("cmdr is already latest", map[string]interface{}{
					"version": core.Version,
				})
				return "", false
			}

			urls := utils.NewSortedHeap(len(release.Assets))
			for _, asset := range release.Assets {
				if asset.BrowserDownloadURL == nil {
					continue
				}

				currentAssetName := asset.GetName()

				if currentAssetName == assetName {
					urls.Add(asset, 0.0)
					break
				}

				score := 0.0
				if strings.Contains(currentAssetName, runtime.GOOS) {
					score += 1
				}
				if strings.Contains(currentAssetName, runtime.GOARCH) {
					score += 1
				}
				if score > 0.0 {
					urls.Add(asset, score)
				}
			}

			item, _ := urls.PopMax()
			if item == nil {
				logger.Warn("no asset found", map[string]interface{}{
					"release": release,
				})
				return "", false
			}

			asset := item.(*github.ReleaseAsset)
			url := asset.GetBrowserDownloadURL()

			logger.Debug("asset url found", map[string]interface{}{
				"url":     url,
				"release": releaseName,
				"asset":   asset.GetName(),
			})

			return url, true
		}

		installCmdr := func(uri string) error {
			downloader := utils.NewProgressBarDownloader(os.Stderr, func(c *getter.Client) error {
				c.Mode = getter.ClientModeFile
				return nil
			})
			err := downloader.Fetch(uri, location)
			if err != nil {
				return errors.Wrapf(err, "failed to download %s", uri)
			}

			err = os.Chmod(location, 0755)
			if err != nil {
				return errors.Wrapf(err, "failed to chmod %s", location)
			}

			return nil
		}

		logger.Debug("searching for release", map[string]interface{}{
			"release":  releaseTag,
			"asset":    assetName,
			"location": location,
		})

		info, err := os.Stat(location)

		if err == nil && info.Mode()&0111 != 0 && utils.WaitProcess(ctx, location, upgradeArgs) == nil {
			return
		}

		url, found := downloadFromGithub()
		if !found {
			return
		}

		err = installCmdr(url)
		if err != nil {
			utils.ExitOnError("download cmdr failed", err)
		}

		utils.ExitOnError("upgrade cmdr failed", utils.WaitProcess(ctx, location, upgradeArgs))
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	cfg := core.GetConfiguration()

	flags := upgradeCmd.Flags()
	flags.StringP("release", "r", "latest", "cmdr release tag name")
	flags.StringP("asset", "a", core.Asset, "cmdr release assert name")
	flags.StringP(
		"location", "l",
		filepath.Join(os.TempDir(), fmt.Sprintf("cmdr_%s_replacement", core.Version)),
		"cmdr binary local location",
	)

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXUpgradeRelease, flags.Lookup("release")),
		cfg.BindPFlag(core.CfgKeyXUpgradeAsset, flags.Lookup("asset")),
		cfg.BindPFlag(core.CfgKeyXUpgradeLocation, flags.Lookup("location")),
	)
}
