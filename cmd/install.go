package cmd

import (
	"path"
	"regexp"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var installCmdFlag struct {
	name     string
	version  string
	location string
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger
		location := installCmdFlag.location
		fs := define.FS
		ctx := cmd.Context()

		httpSchemaRegexp := regexp.MustCompile(`^https?://.*?$`)
		if httpSchemaRegexp.Match([]byte(installCmdFlag.location)) {
			outputDir, err := afero.TempDir(fs, "", "")
			utils.ExitWithError(err, "create temporary dir failed")

			location = path.Join(outputDir, installCmdFlag.name)

			logger.Debug("downloading command", map[string]interface{}{
				"url":    installCmdFlag.location,
				"target": location,
			})
			utils.ExitWithError(
				utils.DownloadToFile(ctx, installCmdFlag.location, location),
				"download command failed",
			)

			logger.Info("command downloaded", map[string]interface{}{
				"url": installCmdFlag.location,
			})
		}

		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)

		logger.Debug("installing command", map[string]interface{}{
			"name":     installCmdFlag.name,
			"version":  installCmdFlag.version,
			"location": location,
		})
		utils.ExitWithError(
			helper.Install(ctx, installCmdFlag.name, installCmdFlag.version, location),
			"install command failed",
		)

		logger.Info("command installed", map[string]interface{}{
			"name":    installCmdFlag.name,
			"version": installCmdFlag.version,
		})
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	flags := installCmd.Flags()
	flags.StringVarP(&installCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&installCmdFlag.version, "version", "v", "", "command version")
	flags.StringVarP(&installCmdFlag.location, "location", "l", "", "command location")

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}
