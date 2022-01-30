package command

import (
	"context"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/operator"
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
	runner := operator.NewOperatorRunner(
		operator.NewDBClientMaker(),
		operator.NewCommandsQuerier(matchers),
		operator.NewCommandSorter(),
		operator.NewCommandHandler("command-filter", handler),
	)
	_ = runner.Run(context.Background())
}

func (f *commandFlagsHelper) registerCommandSuggestions(cmd *cobra.Command, name string, getter func(command *model.Command) string) {
	cmd.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		results := []string{}
		f.queryCommands(f.getMagicMatchers(), func(ctx context.Context, commands []*model.Command) error {
			linq.From(commands).
				GroupBy(func(i interface{}) interface{} {
					return getter(i.(*model.Command))
				}, func(i interface{}) interface{} {
					return i
				}).
				Select(func(i interface{}) interface{} {
					return i.(linq.Group).Key
				}).
				Where(func(i interface{}) bool {
					return strings.HasPrefix(i.(string), toComplete)
				}).
				ToSlice(&results)
			return nil
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
