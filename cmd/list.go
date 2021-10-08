/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model/command"
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
		query := core.NewCommandHelper(client).Query()

		if listCmdFlag.name != "" {
			logger.Debug("filter by name", map[string]interface{}{
				"name": listCmdFlag.name,
			})
			query = query.Where(command.Name(listCmdFlag.name))
		}

		if listCmdFlag.version != "" {
			logger.Debug("filter by version", map[string]interface{}{
				"version": listCmdFlag.version,
			})
			query = query.Where(command.Version(listCmdFlag.version))
		}

		if listCmdFlag.location != "" {
			logger.Debug("filter by location", map[string]interface{}{
				"location": listCmdFlag.location,
			})
			query = query.Where(command.Location(listCmdFlag.location))
		}

		if listCmdFlag.activated {
			logger.Debug("filter by activated", map[string]interface{}{
				"activated": listCmdFlag.activated,
			})
			query = query.Where(command.Activated(listCmdFlag.activated))
		}

		commands, err := query.All(cmd.Context())
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := listCmd.Flags()
	flags.StringVarP(&listCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&listCmdFlag.version, "version", "v", "", "command version")
	flags.StringVarP(&listCmdFlag.location, "location", "l", "", "command location")
}
