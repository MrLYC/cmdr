package core

import (
	"context"
	"fmt"
	"strings"
)

type CommandPrinter struct {
	BaseStep
}

func (p *CommandPrinter) String() string {
	return "printer"
}

func (p *CommandPrinter) Run(ctx context.Context) (context.Context, error) {
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
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
