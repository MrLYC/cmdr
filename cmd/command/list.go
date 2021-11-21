package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var listCmdFlag struct {
	name      string
	version   string
	location  string
	activated bool
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Run: func(cmd *cobra.Command, args []string) {
		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewSimpleCommandsQuerier(
				listCmdFlag.name, listCmdFlag.version, listCmdFlag.location, listCmdFlag.activated,
			),
			core.NewCommandPrinter(),
		)

		utils.ExitWithError(runner.Run(utils.SetIntoContext(cmd.Context(), map[define.ContextKey]interface{}{
			define.ContextKeyName:     simpleCmdFlag.name,
			define.ContextKeyVersion:  simpleCmdFlag.version,
			define.ContextKeyLocation: simpleCmdFlag.location,
		})), "list failed")
	},
}

func init() {
	Cmd.AddCommand(listCmd)

	flags := listCmd.Flags()
	flags.StringVarP(&listCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&listCmdFlag.version, "version", "v", "", "command version")
	flags.StringVarP(&listCmdFlag.location, "location", "l", "", "command location")
}
