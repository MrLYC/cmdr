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
			command, err = helper.GetActivatedCommand(cmd.Context(), execCmdFlag.name)
			utils.CheckError(err)
		} else {
			command, err = helper.GetCommandByNameAndVersion(cmd.Context(), execCmdFlag.name, execCmdFlag.version)
			utils.CheckError(err)
		}

		if command == nil {
			panic(core.ErrCommandNotExists)
		}

		logger.Debug("executing command", map[string]interface{}{
			"target": command.Location,
			"args":   args,
		})
		utils.CheckError(syscall.Exec(command.Location, args, os.Environ()))
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	flags := execCmd.Flags()
	flags.StringVarP(&execCmdFlag.name, "name", "n", "", "command name")
	flags.StringVarP(&execCmdFlag.version, "version", "v", "", "command version")

	execCmd.MarkFlagRequired("name")
}