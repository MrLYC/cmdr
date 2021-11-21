package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

type CommandPrinter struct {
	BaseStep
}

func (p *CommandPrinter) String() string {
	return "printer"
}

func (p *CommandPrinter) Run(ctx context.Context) (context.Context, error) {
	values := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommands)
	if values == nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "commands not found")
	}

	commands, ok := values.([]*model.Command)
	if !ok || len(commands) == 0 {
		return ctx, nil
	}

	for _, command := range commands {
		var parts []string
		if command.Activated {
			parts = append(parts, "*")
		} else {
			parts = append(parts, " ")
		}

		parts = append(parts, fmt.Sprintf("%s@%s", command.Name, command.Version))

		if !command.Managed {
			parts = append(parts, command.Location)
		}

		fmt.Printf("%s\n", strings.Join(parts, " "))
	}

	return ctx, nil
}

func NewCommandPrinter() *CommandPrinter {
	return &CommandPrinter{}
}
