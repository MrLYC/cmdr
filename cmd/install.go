package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/utils"
)

var installCmdFlag struct {
	name     string
	version  string
	location string
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install command into cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)
		utils.CheckError(helper.Install(cmd.Context(), installCmdFlag.name, installCmdFlag.version, installCmdFlag.location))
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	flags := installCmd.Flags()
	flags.StringVarP(&installCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&installCmdFlag.version, "version", "v", "", "command version")
	flags.StringVarP(&installCmdFlag.location, "location", "l", "", "command location")

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}
