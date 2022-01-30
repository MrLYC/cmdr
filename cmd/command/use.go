package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run: func(cmd *cobra.Command, args []string) {
		runner := runner.NewUseRunner(define.Config)

		utils.ExitWithError(runner.Run(cmd.Context()), "activate failed")

		define.Logger.Info("used command", map[string]interface{}{
			"name":    simpleCmdFlag.name,
			"version": simpleCmdFlag.version,
		})
	},
}

func init() {
	Cmd.AddCommand(useCmd)
	cmdFlagsHelper.declareFlagName(useCmd)
	cmdFlagsHelper.declareFlagVersion(useCmd)

	flags := useCmd.Flags()
	cfg := define.Config
	cfg.BindPFlag(runner.CfgKeyCommandUseName, flags.Lookup("name"))
	cfg.BindPFlag(runner.CfgKeyCommandUseVersion, flags.Lookup("version"))

	useCmd.MarkFlagRequired("name")
	useCmd.MarkFlagRequired("version")
}
