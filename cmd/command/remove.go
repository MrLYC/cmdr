package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove command from cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewCommandQuerierByNameAndVersion(
				simpleCmdFlag.name, simpleCmdFlag.version,
			),
			core.NewBinaryRemover(),
			core.NewCommandRemover(),
		)

		utils.ExitWithError(runner.Run(utils.SetIntoContext(cmd.Context(), map[define.ContextKey]interface{}{
			define.ContextKeyName:    simpleCmdFlag.name,
			define.ContextKeyVersion: simpleCmdFlag.version,
		})), "list failed")
	},
}

func init() {
	Cmd.AddCommand(removeCmd)

	flags := removeCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&simpleCmdFlag.version, "version", "v", "", "command version")

	removeCmd.MarkFlagRequired("name")
	removeCmd.MarkFlagRequired("version")
}
