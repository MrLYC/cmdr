package command

import (
	"os"
	"strings"

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
		logger := core.GetLogger()

		commands, err := queryCommands(
			manager,
			cfg.GetBool(core.CfgKeyXCommandListActivate),
			cfg.GetString(core.CfgKeyXCommandListName),
			cfg.GetString(core.CfgKeyXCommandListVersion),
			cfg.GetString(core.CfgKeyXCommandListLocation),
		)
		if err != nil {
			return err
		}

		logger.Debug("query commands", map[string]interface{}{
			"commands": commands,
		})

		fields := cfg.GetStringSlice(core.CfgKeyXCommandListFields)
		rowMaker := func(activateFlag string, name string, version string, location string) []string {
			mappings := map[string]string{
				"activated": activateFlag,
				"name":      name,
				"version":   version,
				"location":  location,
			}

			results := make([]string, 0, len(fields))
			for _, field := range fields {
				result, ok := mappings[strings.ToLower(field)]
				if ok {
					results = append(results, result)
				}
			}

			return results
		}

		tab := table.Table{
			Headers: rowMaker("Activated", "Name", "Version", "Location"),
		}

		for _, cmd := range commands {
			activated := ""
			if cmd.GetActivated() {
				activated = "*"
			}

			tab.Rows = append(tab.Rows, rowMaker(
				activated,
				cmd.GetName(),
				cmd.GetVersion(),
				cmd.GetLocation(),
			))
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
	flags.StringSliceP("fields", "f", []string{"Activated", "Name", "Version", "Location"}, "fields to display")

	cfg := core.GetConfiguration()

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXCommandListName, flags.Lookup("name")),
		cfg.BindPFlag(core.CfgKeyXCommandListVersion, flags.Lookup("version")),
		cfg.BindPFlag(core.CfgKeyXCommandListLocation, flags.Lookup("location")),
		cfg.BindPFlag(core.CfgKeyXCommandListActivate, flags.Lookup("activate")),
		cfg.BindPFlag(core.CfgKeyXCommandListFields, flags.Lookup("fields")),

		utils.NewDefaultCobraCommandCompleteHelper(listCmd).RegisterAll(),
	)
}
