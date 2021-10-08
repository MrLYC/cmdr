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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := installCmd.Flags()
	flags.StringVarP(&installCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&installCmdFlag.version, "version", "v", "", "command version")
	flags.StringVarP(&installCmdFlag.location, "location", "l", "", "command location")

	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("location")
}