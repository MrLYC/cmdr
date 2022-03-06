package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/utils"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "doctor to fix cmdr",
	Run: utils.RunCobraCommandWith(core.CommandProviderDoctor, func(cfg core.Configuration, mgr core.CommandManager) error {
		doctor := manager.NewCommandDoctor(mgr)
		return doctor.Fix()
	}),
}

func init() {
	rootCmd.AddCommand(doctorCmd)

	cfg := core.GetConfiguration()

	flags := doctorCmd.Flags()
	flags.StringP("release", "r", "latest", "cmdr release tag name")
	flags.StringP("asset", "a", core.Asset, "cmdr release assert name")

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXUpgradeRelease, flags.Lookup("release")),
		cfg.BindPFlag(core.CfgKeyXUpgradeAsset, flags.Lookup("asset")),
	)
}
