package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List commands",
	Run: runCommand(func(cfg core.Configuration, manager core.CommandManager) error {
		name := cfg.GetString(core.CfgKeyCommandListName)
		version := cfg.GetString(core.CfgKeyCommandListVersion)
		location := cfg.GetString(core.CfgKeyCommandListLocation)
		activate := cfg.GetBool(core.CfgKeyCommandListActivate)

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
			if command.Activated() {
				parts = append(parts, "*")
			} else {
				parts = append(parts, " ")
			}

			parts = append(parts, command.Name(), command.Version())

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
	cfg.BindPFlag(core.CfgKeyCommandListName, flags.Lookup("name"))
	cfg.BindPFlag(core.CfgKeyCommandListVersion, flags.Lookup("version"))
	cfg.BindPFlag(core.CfgKeyCommandListLocation, flags.Lookup("location"))
	cfg.BindPFlag(core.CfgKeyCommandListActivate, flags.Lookup("activate"))
}
