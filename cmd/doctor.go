package cmd

import (
	"github.com/asdine/storm/v3/q"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

var doctorCmdFlag struct {
	name string
}

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check and fix commands mapping",
	Run: func(cmd *cobra.Command, args []string) {
		client := core.GetClient()
		defer utils.CallClose(client)

		fs := define.FS
		logger := define.Logger
		ctx := cmd.Context()
		filters := make([]q.Matcher, 0)

		logger.Info("rebuild bin dir")
		utils.ExitWithError(
			fs.MkdirAll(core.GetBinDir(), 0755),
			"making dir bin failed",
		)

		if doctorCmdFlag.name != "" {
			logger.Debug("filter by name", map[string]interface{}{
				"name": doctorCmdFlag.name,
			})
			filters = append(filters, q.Eq("Name", doctorCmdFlag.name))
		}

		helper := core.NewCommandHelper(client)

		logger.Debug("quering commands")
		commands, err := helper.GetCommands(ctx, filters...)
		utils.ExitWithError(err, "query command failed")

		for _, command := range commands {
			_, ferr := fs.Stat(command.Location)
			if ferr != nil {
				logger.Debug("deleting command", map[string]interface{}{
					"name":     command.Name,
					"version":  command.Version,
					"location": command.Location,
				})

				err := client.DeleteStruct(command)
				if err != nil {
					logger.Error("remove command failed", map[string]interface{}{
						"name":    command.Name,
						"version": command.Version,
						"error":   err,
					})
				} else {
					logger.Info("command deleted", map[string]interface{}{
						"name":    command.Name,
						"version": command.Version,
					})
				}
			} else if command.Activated {
				logger.Debug("activating command", map[string]interface{}{
					"name":    command.Name,
					"version": command.Version,
				})

				err := helper.Activate(ctx, command.Name, command.Version)
				if err != nil {
					logger.Error("activate command failed", map[string]interface{}{
						"name":    command.Name,
						"version": command.Version,
						"error":   err,
					})
				} else {
					logger.Info("command activated", map[string]interface{}{
						"name":    command.Name,
						"version": command.Version,
					})
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	flags := doctorCmd.Flags()

	flags.StringVarP(&doctorCmdFlag.name, "name", "n", "", "command name")
}
