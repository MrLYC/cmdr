package command

import (
	"context"

	"github.com/asdine/storm/v3/q"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/model"
)

var cmdFlagsHelper commandFlagsHelper

type commandFlagsHelper struct{}

func (f *commandFlagsHelper) declareFlagName(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&simpleCmdFlag.name, "name", "n", "", "command name")
	f.registerCommandSuggestions(cmd, "name", func(command *model.Command) string {
		return command.Name
	})
}

func (f *commandFlagsHelper) declareFlagVersion(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&simpleCmdFlag.version, "version", "v", "", "command version")
	f.registerCommandSuggestions(cmd, "version", func(command *model.Command) string {
		return command.Version
	})
}

func (f *commandFlagsHelper) declareFlagLocation(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&simpleCmdFlag.location, "location", "l", "", "command location")
	f.registerCommandSuggestions(cmd, "location", func(command *model.Command) string {
		return command.Location
	})
}

func (f *commandFlagsHelper) queryCommands(matchers []q.Matcher, handler func(ctx context.Context, commands []*model.Command) error) {
	runner := core.NewStepRunner(
		core.NewDBClientMaker(),
		core.NewCommandsQuerier(matchers),
		core.NewCommandSorter(),
		core.NewCommandHandler("command-filter", handler),
	)
	_ = runner.Run(context.Background())
}

func (f *commandFlagsHelper) registerCommandSuggestions(cmd *cobra.Command, name string, getter func(command *model.Command) string) {
	cmd.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		set := treeset.NewWithStringComparator()
		f.queryCommands(f.getMagicMatchers(), func(ctx context.Context, commands []*model.Command) error {
			for _, command := range commands {
				set.Add(getter(command))
			}
			return nil
		})

		results := make([]string, 0, set.Size())
		set.Each(func(index int, value interface{}) {
			results = append(results, value.(string))
		})

		return results, cobra.ShellCompDirectiveDefault
	})
}

func (f *commandFlagsHelper) getMagicMatchers() []q.Matcher {
	matchers := make([]q.Matcher, 0, 3)

	if simpleCmdFlag.name != "" {
		matchers = append(matchers, q.Eq("Name", simpleCmdFlag.name))
	}
	if simpleCmdFlag.version != "" {
		matchers = append(matchers, q.Eq("Version", simpleCmdFlag.version))
	}
	if simpleCmdFlag.location != "" {
		matchers = append(matchers, q.Eq("Location", simpleCmdFlag.location))
	}

	return matchers
}
