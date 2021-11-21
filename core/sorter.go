package core

import (
	"context"
	"sort"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

type CommandSorter struct {
	BaseStep
}

func (s *CommandSorter) String() string {
	return "command-sorter"
}

func (s *CommandSorter) Run(ctx context.Context) (context.Context, error) {
	values := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommands)
	commands, ok := values.([]*model.Command)
	if !ok || len(commands) == 0 {
		return ctx, nil
	}

	sort.Slice(commands, func(i, j int) bool {
		x := commands[i]
		y := commands[j]

		if x.Activated != y.Activated {
			return !y.Activated
		}

		if x.Name != y.Name {
			return x.Name < y.Name
		}

		return y.Version < x.Version
	})

	return ctx, nil
}

func NewCommandSorter() *CommandSorter {
	return &CommandSorter{}
}
