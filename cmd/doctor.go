package cmd

import (
	"github.com/asdine/storm/v3/q"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
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

		cmdr, err := core.NewCmdr(cfg.GetString(config.CfgKeyCmdrRoot))
		if err != nil {
			utils.ExitWithError(err, "create cmdr failed")
		}
		shimsDir := cmdr.BinaryManager.ShimsManager.Path()
		binDir := cmdr.BinaryManager.BinManager.Path()

		runner := runner.New(
			operator.NewDBMigrator(cmdr),
			operator.NewCommandsQuerier([]q.Matcher{q.Eq("Activated", true)}),
			operator.NewBrokenCommandsFixer(cmdr),
			operator.NewDirectoryRemover(map[string]string{
				"bin": binDir,
			}),
			operator.NewDirectoryMaker(map[string]string{
				"shims": shimsDir,
				"bin":   binDir,
			}),
			operator.NewBinariesActivator(cmdr),
		)

		utils.ExitWithError(runner.Run(cmd.Context()), "doctor failed")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	flags := doctorCmd.Flags()

	flags.StringVarP(&doctorCmdFlag.name, "name", "n", "", "command name")
}
