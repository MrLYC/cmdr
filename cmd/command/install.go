package command

import (
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
		shimsDir := core.GetShimsDir()

		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewDownloader(),
			core.NewCommandDefiner(shimsDir),
			core.NewBinariesInstaller(shimsDir),
		)

		utils.ExitWithError(runner.Run(utils.SetIntoContext(cmd.Context(), map[define.ContextKey]interface{}{
			define.ContextKeyName:           simpleCmdFlag.name,
			define.ContextKeyVersion:        simpleCmdFlag.version,
			define.ContextKeyLocation:       simpleCmdFlag.location,
			define.ContextKeyCommandManaged: true,
		})), "install failed")

		define.Logger.Info("installed command", map[string]interface{}{
			"name":     simpleCmdFlag.name,
			"version":  simpleCmdFlag.version,
			"location": simpleCmdFlag.location,
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
