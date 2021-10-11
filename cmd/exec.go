//+build !windows

package cmd

import (
	"os"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

var execCmdFlag struct {
	name    string
	version string
}

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command",
	Run: func(cmd *cobra.Command, args []string) {
		logger := define.Logger
		client := core.GetClient()
		defer utils.CallClose(client)

		helper := core.NewCommandHelper(client)
		var (
			command *model.Command
			err     error
		)

		if execCmdFlag.version == "" {
			logger.Debug("looking up the activated command", map[string]interface{}{
				"name": execCmdFlag.name,
			})

			command, err = helper.GetActivatedCommand(cmd.Context(), execCmdFlag.name)
			utils.ExitWithError(err, "lookup activated command failed")
		} else {
			logger.Debug("looking up the command", map[string]interface{}{
				"name":    execCmdFlag.name,
				"version": execCmdFlag.version,
			})

			command, err = helper.GetCommandByNameAndVersion(cmd.Context(), execCmdFlag.name, execCmdFlag.version)
			utils.ExitWithError(err, "lookup specified command failed")
		}

		if command == nil {
			logger.Warn("command not found", map[string]interface{}{
				"name":    execCmdFlag.name,
				"version": execCmdFlag.version,
			})
			exitCode = -2
			return
		}

		logger.Debug("executing", map[string]interface{}{
			"name":    execCmdFlag.name,
			"version": execCmdFlag.version,
			"target":  command.Location,
			"args":    args,
		})
		utils.ExitWithError(
			syscall.Exec(command.Location, args, os.Environ()),
			"execute command failed",
		)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	flags := execCmd.Flags()
	flags.StringVarP(&execCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&execCmdFlag.version, "version", "v", "", "command version")

	execCmd.MarkFlagRequired("name")
}
