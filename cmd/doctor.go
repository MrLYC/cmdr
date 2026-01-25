package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/utils"
)

var (
	doctorDryRun   bool
	doctorNoBackup bool
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "doctor to fix cmdr",
	Run: utils.RunCobraCommandWith(core.CommandProviderDoctor, func(cfg core.Configuration, mgr core.CommandManager) error {
		rootDir := cfg.GetString(core.CfgKeyCmdrRootDir)
		doctor := manager.NewCommandDoctor(mgr, rootDir)
		return doctor.FixWithOptions(doctorDryRun, !doctorNoBackup)
	}),
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVar(&doctorDryRun, "dry-run", false, "show what would be done without making any changes")
	doctorCmd.Flags().BoolVar(&doctorNoBackup, "no-backup", false, "skip backup before making changes")
}
