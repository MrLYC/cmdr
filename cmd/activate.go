package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/utils"
)

var activateCmdFlag struct {
	name    string
	version string
}

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate a command",
	Run: func(cmd *cobra.Command, args []string) {
		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)
		utils.CheckError(helper.Activate(cmd.Context(), activateCmdFlag.name, activateCmdFlag.version))
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)

	flags := activateCmd.Flags()
	flags.StringVarP(&activateCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&activateCmdFlag.version, "version", "v", "", "command version")

	activateCmd.MarkFlagRequired("name")
	activateCmd.MarkFlagRequired("version")
}
