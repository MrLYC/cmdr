package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/utils"
)

var (
	doctorDryRun bool
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "doctor to fix cmdr",
	Run: utils.RunCobraCommandWith(core.CommandProviderDoctor, func(cfg core.Configuration, mgr core.CommandManager) error {
		doctor := manager.NewCommandDoctor(mgr)
		return doctor.Fix(doctorDryRun)
	}),
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVar(&doctorDryRun, "dry-run", false, "show what would be done without making any changes")
}
