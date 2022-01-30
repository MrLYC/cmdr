package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run:   executeRunner(runner.NewUseRunner),
}

func init() {
	Cmd.AddCommand(useCmd)
	flags := useCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")

	cfg := define.Config
	cfg.BindPFlag(runner.CfgKeyCommandUseName, flags.Lookup("name"))
	cfg.BindPFlag(runner.CfgKeyCommandUseVersion, flags.Lookup("version"))

	useCmd.MarkFlagRequired("name")
	useCmd.MarkFlagRequired("version")
}
