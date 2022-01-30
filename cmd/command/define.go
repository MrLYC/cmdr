package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
)

// defineCmd represents the define command
var defineCmd = &cobra.Command{
	Use:   "define",
	Short: "Define command into cmdr",
	Run:   executeRunner(runner.NewDefineRunner),
}

func init() {
	Cmd.AddCommand(defineCmd)
	flags := defineCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := define.Config
	cfg.BindPFlag(runner.CfgKeyCommandDefineName, flags.Lookup("name"))
	cfg.BindPFlag(runner.CfgKeyCommandDefineVersion, flags.Lookup("version"))
	cfg.BindPFlag(runner.CfgKeyCommandDefineLocation, flags.Lookup("location"))

	defineCmd.MarkFlagRequired("name")
	defineCmd.MarkFlagRequired("version")
	defineCmd.MarkFlagRequired("location")
}
