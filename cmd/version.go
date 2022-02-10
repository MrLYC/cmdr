package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

var versionCmdFlag struct {
	all bool
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print cmdr version",
	Run: func(cmd *cobra.Command, args []string) {
		if versionCmdFlag.all {
			fmt.Printf(
				"Author: %s\nVersion: %s\nCommit: %s\nDate: %s\nAsset: %s\n",
				core.Author,
				core.Version,
				core.Commit,
				core.BuildDate,
				core.Asset,
			)

		} else {
			fmt.Println(core.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	flags := versionCmd.Flags()

	flags.BoolVarP(&versionCmdFlag.all, "all", "a", false, "print all infomation")
}
