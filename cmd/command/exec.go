//+build !windows

package command

import (
	"os"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

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

		if simpleCmdFlag.version == "" {
			logger.Debug("looking up the activated command", map[string]interface{}{
				"name": simpleCmdFlag.name,
			})

			command, err = helper.GetActivatedCommand(cmd.Context(), simpleCmdFlag.name)
			utils.ExitWithError(err, "lookup activated command failed")
		} else {
			logger.Debug("looking up the command", map[string]interface{}{
				"name":    simpleCmdFlag.name,
				"version": simpleCmdFlag.version,
			})

			command, err = helper.GetCommandByNameAndVersion(cmd.Context(), simpleCmdFlag.name, simpleCmdFlag.version)
			utils.ExitWithError(err, "lookup specified command failed")
		}

		if command == nil {
			logger.Warn("command not found", map[string]interface{}{
				"name":    simpleCmdFlag.name,
				"version": simpleCmdFlag.version,
			})
			return
		}

		logger.Debug("executing", map[string]interface{}{
			"name":    simpleCmdFlag.name,
			"version": simpleCmdFlag.version,
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
	Cmd.AddCommand(execCmd)

	flags := execCmd.Flags()
	flags.StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&simpleCmdFlag.version, "version", "v", "", "command version")

	execCmd.MarkFlagRequired("name")
}
