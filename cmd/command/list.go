package command

import (
	"os"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
	"github.com/tomlazar/table"
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

		if name != "" {
			query.WithName(name)
		}

		if version != "" {
			query.WithVersion(version)
		}

		if location != "" {
			query.WithLocation(location)
		}

		if activate {
			query.WithActivated(activate)
		}

		commands, err := query.All()
		if err != nil {
			return err
		}

		tab := table.Table{
			Headers: []string{"Activated", "Name", "Version", "Location"},
		}

		utils.SortCommands(commands)
		for _, cmd := range commands {
			activated := ""
			if cmd.GetActivated() {
				activated = "*"
			}

			tab.Rows = append(tab.Rows, []string{
				activated,
				cmd.GetName(),
				cmd.GetVersion(),
				cmd.GetLocation(),
			})
		}

		return tab.WriteTable(os.Stdout, &table.Config{
			Color:           true,
			AlternateColors: true,
			TitleColorCode:  ansi.ColorCode("white+buf"),
		})
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

		utils.NewDefaultCobraCommandCompleteHelper(listCmd).RegisterAll(),
	)
}
