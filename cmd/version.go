package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
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
				"Author: %s\nVersion: %s\nCommit: %s\nDate: %s\nOS: %s\nArch: %s\n",
				define.Author,
				define.Version,
				define.Commit,
				define.BuildDate,
				runtime.GOOS,
				runtime.GOARCH,
			)

		} else {
			fmt.Println(define.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	flags := versionCmd.Flags()

	flags.BoolVarP(&versionCmdFlag.all, "all", "a", false, "print all infomation")
}
