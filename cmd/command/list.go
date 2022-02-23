package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		name := cfg.GetString(core.CfgKeyXCommandListName)
		version := cfg.GetString(core.CfgKeyXCommandListVersion)
		location := cfg.GetString(core.CfgKeyXCommandListLocation)
		activate := cfg.GetBool(core.CfgKeyXCommandListActivate)

		query, err := manager.Query()
		if err != nil {
			return err
		}

		switch {
		case name != "":
			query.WithName(name)
		case version != "":
			query.WithVersion(version)
		case location != "":
			query.WithLocation(location)
		case activate:
			query.WithActivated(activate)
		}

		commands, err := query.All()
		if err != nil {
			return err
		}

		for _, command := range commands {
			var parts []string
			if command.GetActivated() {
				parts = append(parts, "*")
			} else {
				parts = append(parts, " ")
			}

			parts = append(parts, command.GetName(), command.GetVersion())

			_, _ = fmt.Fprintf(os.Stdout, "%s\n", strings.Join(parts, " "))
		}

		return nil
	}),
}

func init() {
	Cmd.AddCommand(listCmd)
	flags := listCmd.Flags()
	flags.StringP("name", "n", "", "command name")
	flags.StringP("version", "v", "", "command version")
	flags.StringP("location", "l", "", "command location")
	flags.BoolP("activate", "a", false, "activate command")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXCommandListName, flags.Lookup("name")),
		cfg.BindPFlag(core.CfgKeyXCommandListVersion, flags.Lookup("version")),
		cfg.BindPFlag(core.CfgKeyXCommandListLocation, flags.Lookup("location")),
		cfg.BindPFlag(core.CfgKeyXCommandListActivate, flags.Lookup("activate")),
	)
}
