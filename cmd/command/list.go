package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Run:   executeRunner(runner.NewListRunner),
}

func init() {
	Cmd.AddCommand(listCmd)
	flags := listCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")
	flags.BoolP("activate", "a", false, "activate command")

	cfg := define.Config
	cfg.BindPFlag(runner.CfgKeyCommandListName, flags.Lookup("name"))
	cfg.BindPFlag(runner.CfgKeyCommandListVersion, flags.Lookup("version"))
	cfg.BindPFlag(runner.CfgKeyCommandListLocation, flags.Lookup("location"))
	cfg.BindPFlag(runner.CfgKeyCommandListActivate, flags.Lookup("activate"))
}
