package command

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Activate a command",
	Run: func(cmd *cobra.Command, args []string) {
		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewCommandQuerierByNameAndVersion(simpleCmdFlag.name, simpleCmdFlag.version),
			core.NewBinaryActivator(),
			core.NewCommandDeactivator(),
			core.NewCommandActivator(),
		)

		utils.ExitWithError(runner.Run(utils.SetIntoContext(cmd.Context(), map[define.ContextKey]interface{}{
			define.ContextKeyName:     simpleCmdFlag.name,
			define.ContextKeyVersion:  simpleCmdFlag.version,
			define.ContextKeyLocation: simpleCmdFlag.location,
		})), "activate failed")
	},
}

func init() {
	Cmd.AddCommand(useCmd)

	flags := useCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&simpleCmdFlag.version, "version", "v", "", "command version")

	useCmd.MarkFlagRequired("name")
	useCmd.MarkFlagRequired("version")
}
