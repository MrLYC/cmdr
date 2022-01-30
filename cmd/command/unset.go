package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Deactivate a command",
	Run: func(cmd *cobra.Command, args []string) {
		runner := runner.NewUnsetRunner(define.Config)

		utils.ExitWithError(runner.Run(cmd.Context()), "deactivate failed")

		define.Logger.Info("unset command", map[string]interface{}{
			"name": simpleCmdFlag.name,
		})
	},
}

func init() {
	Cmd.AddCommand(unsetCmd)
	cmdFlagsHelper.declareFlagName(unsetCmd)

	flags := unsetCmd.Flags()
	cfg := define.Config
	cfg.BindPFlag(runner.CfgKeyCommandUnsetName, flags.Lookup("name"))

	unsetCmd.MarkFlagRequired("name")
}
