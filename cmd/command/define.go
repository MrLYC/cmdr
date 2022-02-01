package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
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
	flags.StringP("location", "l", "", "command location")

	cfg := config.Global
	cfg.BindPFlag(config.CfgKeyCommandDefineName, flags.Lookup("name"))
	cfg.BindPFlag(config.CfgKeyCommandDefineVersion, flags.Lookup("version"))
	cfg.BindPFlag(config.CfgKeyCommandDefineLocation, flags.Lookup("location"))

	defineCmd.MarkFlagRequired("name")
	defineCmd.MarkFlagRequired("version")
	defineCmd.MarkFlagRequired("location")
}
