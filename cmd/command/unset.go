package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/runner"
)

// unsetCmd represents the unset command
var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Deactivate a command",
	Run:   executeRunner(runner.NewUnsetRunner),
}

func init() {
	Cmd.AddCommand(unsetCmd)
	flags := unsetCmd.Flags()
	flags.StringP("name", "n", "", "command name")

	cfg := config.Global
	cfg.BindPFlag(config.CfgKeyCommandUnsetName, flags.Lookup("name"))

	unsetCmd.MarkFlagRequired("name")
}
