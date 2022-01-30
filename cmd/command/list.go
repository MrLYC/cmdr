package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/utils"
)

var listCmdFlag struct {
	activated bool
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Run: func(cmd *cobra.Command, args []string) {
		runner := runner.NewListRunner(define.Config)
		utils.ExitWithError(runner.Run(cmd.Context()), "list failed")
	},
}

func init() {
	Cmd.AddCommand(listCmd)

	cmdFlagsHelper.declareFlagName(listCmd)
	cmdFlagsHelper.declareFlagVersion(listCmd)
	cmdFlagsHelper.declareFlagLocation(listCmd)

	flags := installCmd.Flags()
	cfg := define.Config
	cfg.BindPFlag(runner.CfgKeyCommandListName, flags.Lookup("name"))
	cfg.BindPFlag(runner.CfgKeyCommandListVersion, flags.Lookup("version"))
	cfg.BindPFlag(runner.CfgKeyCommandListLocation, flags.Lookup("location"))

}
