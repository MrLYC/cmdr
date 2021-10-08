package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model/command"
	"github.com/mrlyc/cmdr/model/predicate"
	"github.com/mrlyc/cmdr/model/schema"
	"github.com/mrlyc/cmdr/utils"
)

var listCmdFlag struct {
	name      string
	version   string
	location  string
	activated bool
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all commands",
	Run: func(cmd *cobra.Command, args []string) {
		client := core.GetClient()
		defer utils.CallClose(client)

		logger := define.Logger
		filters := make([]predicate.Command, 0)

		if listCmdFlag.name != "" {
			logger.Debug("filter by name", map[string]interface{}{
				"name": listCmdFlag.name,
			})
			filters = append(filters, command.Name(listCmdFlag.name))
		}

		if listCmdFlag.version != "" {
			logger.Debug("filter by version", map[string]interface{}{
				"version": listCmdFlag.version,
			})
			filters = append(filters, command.Version(listCmdFlag.version))
		}

		if listCmdFlag.location != "" {
			logger.Debug("filter by location", map[string]interface{}{
				"location": listCmdFlag.location,
			})
			filters = append(filters, command.Location(listCmdFlag.location))
		}

		if listCmdFlag.activated {
			logger.Debug("filter by activated", map[string]interface{}{
				"activated": listCmdFlag.activated,
			})
			filters = append(filters, command.Activated(listCmdFlag.activated))
		}

		commands, err := core.NewCommandHelper(client).GetCommands(cmd.Context(), filters...)
		utils.CheckError(err)

		table := core.NewModleTablePrinter(schema.Command{}, os.Stdout)
		for _, command := range commands {
			utils.CheckError(table.Append(command))
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	flags := listCmd.Flags()
	flags.StringVarP(&listCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&listCmdFlag.version, "version", "v", "", "command version")
	flags.StringVarP(&listCmdFlag.location, "location", "l", "", "command location")
}
