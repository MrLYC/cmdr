package cmd

import (
	"github.com/asdine/storm/v3/q"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/runner"
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
		cfg := config.GetGlobalConfiguration()
		binDir := config.GetBinDir(cfg)

		runner := runner.New(
			operator.NewDBClientMaker(),
			operator.NewDBMigrator(new(model.Command)),
			operator.NewCommandsQuerier([]q.Matcher{q.Eq("Activated", true)}),
			operator.NewBrokenCommandsFixer(),
			operator.NewDirectoryRemover(map[string]string{
				"bin": binDir,
			}),
			operator.NewDirectoryMaker(map[string]string{
				"shims": config.GetShimsDir(cfg),
				"bin":   binDir,
			}),
			operator.NewBinariesInstaller(),
			operator.NewBinariesActivator(),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "doctor failed")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	flags := doctorCmd.Flags()

	flags.StringVarP(&doctorCmdFlag.name, "name", "n", "", "command name")
}
