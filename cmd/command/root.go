package command

import "github.com/spf13/cobra"

var simpleCmdFlag struct {
	name     string
	version  string
	location string
}

var Cmd = &cobra.Command{
	Use:   "command",
	Short: "Manage commands",
}
