package cmd

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var upgradeCmdFlag struct {
	release string
}

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			logger = define.Logger
			ctx    = cmd.Context()
			fs     = define.FS
		)

		logger.Debug("searching cmdr release", map[string]interface{}{
			"release": upgradeCmdFlag.release,
		})
		release, err := utils.GetCMDRRelease(ctx, upgradeCmdFlag.release)
		utils.ExitWithError(err, "search cmdr releases failed")

		outputDir, err := afero.TempDir(fs, "", "")
		utils.ExitWithError(err, "create temporary dir failed")

		version := strings.TrimPrefix(*release.TagName, "v")
		assetName := fmt.Sprintf(
			"%s_%s_%s_%s.tar.gz",
			define.Name,
			version,
			runtime.GOOS,
			runtime.GOARCH,
		)
		logger.Info("cmdr release found", map[string]interface{}{
			"version": version,
		})

		target := path.Join(outputDir, assetName)
		logger.Debug("searching cmdr asset", map[string]interface{}{
			"release": upgradeCmdFlag.release,
			"asset":   assetName,
			"version": version,
			"target":  target,
		})
		utils.ExitWithError(
			utils.DownloadReleaseAssetByName(ctx, release, assetName, target),
			"download asset %s failed", assetName,
		)
		logger.Info("asset downloaded", map[string]interface{}{
			"asset":   assetName,
			"version": version,
		})

		logger.Debug("extraing asset", map[string]interface{}{
			"target": target,
			"output": outputDir,
		})
		utils.ExitWithError(
			utils.ExtraTGZ(target, outputDir),
			"extra asset failed",
		)
		logger.Info("asset unpacked")

		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)
		cmdrPath := path.Join(outputDir, define.Name)

		logger.Debug("upgrading cmdr", map[string]interface{}{
			"version": version,
			"path":    cmdrPath,
		})
		installed, err := helper.Upgrade(ctx, version, cmdrPath)
		utils.ExitWithError(err, "upgrade cmdr failed")
		logger.Info("cmdr upgraded", map[string]interface{}{
			"version": version,
		})

		if installed {
			logger.Info("cmdr already installed")
		} else {
			logger.Info("cmdr installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	flags := upgradeCmd.Flags()
	flags.StringVarP(&upgradeCmdFlag.release, "release", "r", "latest", "cmdr release tag name")
}
