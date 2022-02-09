package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/cmdr"
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
				cmdr.Author,
				cmdr.Version,
				cmdr.Commit,
				cmdr.BuildDate,
				cmdr.Asset,
			)

		} else {
			fmt.Println(cmdr.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	flags := versionCmd.Flags()

	flags.BoolVarP(&versionCmdFlag.all, "all", "a", false, "print all infomation")
}
