package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Deactivate a command",
	Run: runCommand(func(cfg cmdr.Configuration, manager cmdr.CommandManager) error {
		return manager.Deactivate(cfg.GetString(cmdr.CfgKeyCommandUnsetName))
	}),
}

func init() {
	Cmd.AddCommand(unsetCmd)
	flags := unsetCmd.Flags()
	flags.StringP("name", "n", "", "command name")

	cfg := cmdr.GetConfiguration()
	cfg.BindPFlag(cmdr.CfgKeyCommandUnsetName, flags.Lookup("name"))

	unsetCmd.MarkFlagRequired("name")
}
