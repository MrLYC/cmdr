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

		downloadFromGithub := func() string {
			release, err := utils.GetCMDRRelease(ctx, releaseTag)
			if err != nil {
				utils.ExitOnError("search cmdr release failed", err)
			}

			urls := utils.NewSortedHeap(len(release.Assets))
			for _, asset := range release.Assets {
				if asset.BrowserDownloadURL == nil {
					continue
				}

				if *asset.Name == assetName {
					urls.Add(asset, 0.0)
					break
				}

				score := 0.0
				if strings.Contains(*asset.Name, runtime.GOOS) {
					score += 1
				}
				if strings.Contains(*asset.Name, runtime.GOARCH) {
					score += 1
				}
				urls.Add(asset, score)
			}

			item, _ := urls.PopMax()
			if item == nil {
				utils.ExitOnError("no asset found")
			}

			assert := item.(*github.ReleaseAsset)
			url := *assert.BrowserDownloadURL

			logger.Debug("asset url found", map[string]interface{}{
				"url":     url,
				"release": releaseTag,
				"asset":   *assert.Name,
			})

			return url
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

		if err != nil || info.Mode()&0111 == 0 {
			err := installCmdr(downloadFromGithub())
			if err != nil {
				utils.ExitOnError("download cmdr failed", err)
			}
		}

		utils.ExitOnError(
			"Reinitialize cmdr",
			utils.WaitProcess(ctx, location, cfg.GetStringSlice(core.CfgKeyXUpgradeArgs)),
		)
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
