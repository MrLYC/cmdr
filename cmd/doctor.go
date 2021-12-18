package cmd

import (
	"github.com/asdine/storm/v3/q"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
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
		shimsDir := core.GetShimsDir()

		runner := core.NewStepRunner(
			core.NewDBClientMaker(),
			core.NewDBMigrator(new(model.Command)),
			core.NewCommandsQuerier([]q.Matcher{q.Eq("Activated", true)}),
			core.NewBrokenCommandsFixer(shimsDir),
			core.NewDirectoryRemover(map[string]string{
				"bin": binDir,
			}),
			core.NewDirectoryMaker(map[string]string{
				"shims": core.GetShimsDir(),
				"bin":   binDir,
			}),
			core.NewBinariesInstaller(shimsDir),
			core.NewBinariesActivator(binDir),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "doctor failed")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	flags := doctorCmd.Flags()

	flags.StringVarP(&doctorCmdFlag.name, "name", "n", "", "command name")
}
