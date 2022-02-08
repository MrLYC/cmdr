package operator

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type CommandPrinter struct {
	BaseOperator
	writer io.Writer
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

		parts = append(parts, command.Name, command.Version)

		fmt.Fprintf(p.writer, "%s\n", strings.Join(parts, " "))
	}

	return ctx, nil
}

func NewCommandPrinter(writer io.Writer) *CommandPrinter {
	return &CommandPrinter{
		writer: writer,
	}
}
