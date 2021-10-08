package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/utils"
)

var removeCmdFlag struct {
	name    string
	version string
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove command from cmdr",
	Run: func(cmd *cobra.Command, args []string) {
		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)
		utils.CheckError(helper.Remove(cmd.Context(), removeCmdFlag.name, removeCmdFlag.version))
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	flags := removeCmd.Flags()
	flags.StringVarP(&removeCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&removeCmdFlag.version, "version", "v", "", "command version")

	removeCmd.MarkFlagRequired("name")
	removeCmd.MarkFlagRequired("version")
}
