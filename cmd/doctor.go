package cmd

import (
	"github.com/asdine/storm/v3/q"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

var doctorCmdFlag struct {
	name string
}

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check and fix cmdr environment",
	Run: func(cmd *cobra.Command, args []string) {
		binDir := core.GetBinDir()

		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewDBMigrator(new(model.Command)),
			core.NewCommandsQuerier([]q.Matcher{q.Eq("Activated", true)}),
			core.NewBrokenCommandsFixer(),
			core.NewDirectoryRemover(map[string]string{
				"bin": binDir,
			}),
			core.NewDirectoryMaker(map[string]string{
				"shims": core.GetShimsDir(),
				"bin":   binDir,
			}),
			core.NewBinariesInstaller(),
			core.NewBinariesActivator(),
		)

		utils.ExitWithError(runner.Run(utils.SetIntoContext(cmd.Context(), map[define.ContextKey]interface{}{
			define.ContextKeyName:           define.Name,
			define.ContextKeyCommandManaged: true,
		})), "doctor failed")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	flags := doctorCmd.Flags()

	flags.StringVarP(&doctorCmdFlag.name, "name", "n", "", "command name")
}
