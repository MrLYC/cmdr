package command

import (
	"path"
	"regexp"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger
		location := simpleCmdFlag.location
		fs := define.FS
		ctx := cmd.Context()

		httpSchemaRegexp := regexp.MustCompile(`^https?://.*?$`)
		if httpSchemaRegexp.Match([]byte(simpleCmdFlag.location)) {
			outputDir, err := afero.TempDir(fs, "", "")
			utils.ExitWithError(err, "create temporary dir failed")

			location = path.Join(outputDir, simpleCmdFlag.name)

			logger.Debug("downloading command", map[string]interface{}{
				"url":    simpleCmdFlag.location,
				"target": location,
			})
			utils.ExitWithError(
				utils.DownloadToFile(ctx, simpleCmdFlag.location, location),
				"download command failed",
			)

			logger.Info("command downloaded", map[string]interface{}{
				"url": simpleCmdFlag.location,
			})
		}

		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)

		logger.Debug("installing command", map[string]interface{}{
			"name":     simpleCmdFlag.name,
			"version":  simpleCmdFlag.version,
			"location": location,
		})
		utils.ExitWithError(
			helper.Install(ctx, simpleCmdFlag.name, simpleCmdFlag.version, location),
			"install command failed",
		)

		logger.Info("command installed", map[string]interface{}{
			"name":    simpleCmdFlag.name,
			"version": simpleCmdFlag.version,
		})
	},
}

func init() {
	Cmd.AddCommand(installCmd)

	flags := installCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&simpleCmdFlag.version, "version", "v", "", "command version")
	flags.StringVarP(&simpleCmdFlag.location, "location", "l", "", "command location")

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}
