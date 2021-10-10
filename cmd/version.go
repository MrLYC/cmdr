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
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
)

var versionCmdFlag struct {
	all bool
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print cmdr version",
	Run: func(cmd *cobra.Command, args []string) {
		if versionCmdFlag.all {
			fmt.Printf(
				"Author: %s\nVersion: %s\nCommit: %s\nDate: %s\nOS: %s\nArch: %s\n",
				define.Author,
				define.Version,
				define.Commit,
				define.BuildDate,
				runtime.GOOS,
				runtime.GOARCH,
			)

		} else {
			fmt.Println(define.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionCmdFlag.all, "all", "a", false, "print all infomation")
}
