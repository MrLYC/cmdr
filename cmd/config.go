package cmd

import "github.com/mrlyc/cmdr/cmd/config"

func init() {
	rootCmd.AddCommand(config.Cmd)
}
