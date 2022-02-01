package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
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

	cfg := config.Global
	cfg.BindPFlag(config.CfgKeyCommandListName, flags.Lookup("name"))
	cfg.BindPFlag(config.CfgKeyCommandListVersion, flags.Lookup("version"))
	cfg.BindPFlag(config.CfgKeyCommandListLocation, flags.Lookup("location"))
	cfg.BindPFlag(config.CfgKeyCommandListActivate, flags.Lookup("activate"))
}
