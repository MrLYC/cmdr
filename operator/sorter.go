package operator

import (
	"context"
	"sort"
)

type CommandSorter struct {
	BaseOperator
}

func (s *CommandSorter) String() string {
	return "command-sorter"
}

func (s *CommandSorter) Run(ctx context.Context) (context.Context, error) {
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
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
