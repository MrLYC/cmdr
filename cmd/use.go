package cmd

import "github.com/mrlyc/cmdr/cmd/command"

func init() {
	rootCmd.AddCommand(command.UseCmd)
}
