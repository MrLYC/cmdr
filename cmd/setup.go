package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
	"github.com/mrlyc/cmdr/cmdr/utils"
)

var setupCmdFlag struct {
	skipInstall bool
	skipProfile bool
	upgrade     bool
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmdr.GetConfiguration()
		manager, err := cmdr.NewCommandManager(cmdr.CommandProviderSimple, cfg)
		utils.ExitWithError(err, "Failed to create command manager")

		utils.ExitWithError(manager.Init(), "Failed to init command manager")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	flags := setupCmd.Flags()
	flags.BoolVar(&setupCmdFlag.skipInstall, "skip-install", false, "do not install cmdr")
	flags.BoolVar(&setupCmdFlag.skipProfile, "skip-profile", false, "do not write profile")
	flags.BoolVar(&setupCmdFlag.upgrade, "upgrade", false, "for upgrade setup")
}
