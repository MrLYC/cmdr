package cmd

import (
	"os"

	"github.com/asdine/storm/v3/q"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
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
		filters := make([]q.Matcher, 0)

		if listCmdFlag.name != "" {
			logger.Debug("filter by name", map[string]interface{}{
				"name": listCmdFlag.name,
			})
			filters = append(filters, q.Eq("Name", listCmdFlag.name))
		}

		if listCmdFlag.version != "" {
			logger.Debug("filter by version", map[string]interface{}{
				"version": listCmdFlag.version,
			})
			filters = append(filters, q.Eq("Version", listCmdFlag.version))
		}

		if listCmdFlag.location != "" {
			logger.Debug("filter by location", map[string]interface{}{
				"location": listCmdFlag.location,
			})
			filters = append(filters, q.Eq("Location", listCmdFlag.location))
		}

		if listCmdFlag.activated {
			logger.Debug("filter by activated", map[string]interface{}{
				"activated": listCmdFlag.activated,
			})
			filters = append(filters, q.Eq("Activated", listCmdFlag.activated))
		}

		commands, err := core.NewCommandHelper(client).GetCommands(cmd.Context(), filters...)
		utils.CheckError(err)

		table := core.NewModleTablePrinter(os.Stdout)
		table.SetHeader([]string{
			"Activated", "Name", "Version", "Location", "Managed",
		})
		for _, command := range commands {
			table.Append([]string{
				cast.ToString(command.Activated),
				command.Name,
				command.Version,
				command.Location,
				cast.ToString(command.Managed),
			})
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
