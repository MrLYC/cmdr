package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

var installCmdFlag struct {
	activate bool
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		binDir := operator.GetBinDir()
		shimsDir := operator.GetShimsDir()

		runner := operator.NewOperatorRunner(
			operator.NewDBClientMaker(),
			operator.NewCommandDefiner(shimsDir, simpleCmdFlag.name, simpleCmdFlag.version, simpleCmdFlag.location, true),
			operator.NewDownloader(),
			operator.NewBinariesInstaller(shimsDir),
		)

		if installCmdFlag.activate {
			runner.Add(
				operator.NewCommandDeactivator(),
				operator.NewBinariesActivator(binDir, shimsDir),
				operator.NewCommandActivator(),
			)
		}

		utils.ExitWithError(runner.Run(cmd.Context()), "install failed")

		define.Logger.Info("installed command", map[string]interface{}{
			"name":     simpleCmdFlag.name,
			"version":  simpleCmdFlag.version,
			"location": simpleCmdFlag.location,
		})
	},
}

func init() {
	Cmd.AddCommand(installCmd)
	cmdFlagsHelper.declareFlagName(installCmd)
	cmdFlagsHelper.declareFlagVersion(installCmd)

	flags := installCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.location, "location", "l", "", "command location")
	flags.BoolVarP(&installCmdFlag.activate, "activate", "a", false, "activate command")

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}
