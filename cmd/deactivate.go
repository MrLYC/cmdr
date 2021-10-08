package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/utils"
)

var deactivateCmdFlag struct {
	name string
}

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Activate a command",
	Run: func(cmd *cobra.Command, args []string) {
		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)
		utils.CheckError(helper.Deactivate(cmd.Context(), deactivateCmdFlag.name))
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)

	flags := deactivateCmd.Flags()
	flags.StringVarP(&deactivateCmdFlag.name, "name", "n", "", "command name")

	deactivateCmd.MarkFlagRequired("name")
}
