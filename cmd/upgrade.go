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
		utils.CheckError(err)

		outputDir, err := afero.TempDir(fs, "", "")
		utils.CheckError(err)

		version := strings.TrimPrefix(*release.TagName, "v")
		assetName := fmt.Sprintf(
			"%s_%s_%s_%s.tar.gz",
			define.Name,
			version,
			runtime.GOOS,
			runtime.GOARCH,
		)
		target := path.Join(outputDir, assetName)
		logger.Info("searching cmdr assert", map[string]interface{}{
			"release": upgradeCmdFlag.release,
			"assert":  assetName,
			"version": version,
			"target":  target,
		})

		utils.CheckError(utils.DownloadReleaseAssertByName(ctx, release, assetName, target))
		utils.CheckError(utils.ExtraTGZ(target, outputDir))

		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)
		installed, err := helper.Upgrade(ctx, version, path.Join(outputDir, define.Name))
		utils.CheckError(err)

		if installed {
			logger.Info("cmdr already installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	flags := upgradeCmd.Flags()
	flags.StringVarP(&upgradeCmdFlag.release, "release", "r", "latest", "cmdr release tag name")
}
