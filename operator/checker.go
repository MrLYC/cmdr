package operator

import (
	"context"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type BinariesChecker struct {
	BaseOperator
}

func (c *BinariesChecker) String() string {
	return "binaries-checker"
}

func (c *BinariesChecker) Run(ctx context.Context) (context.Context, error) {
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get commands from context failed")
	}

	var errs error
	for _, command := range commands {
		_, err := os.Stat(command.Location)

		if err != nil {
			errs = multierror.Append(errs, errors.WithMessage(err, command.Location))
		}
	}

	return ctx, errs
}

func NewBinariesChecker() *BinariesChecker {
	return &BinariesChecker{}
}

type CommandsChecker struct {
	BaseOperator
}

func (c *CommandsChecker) String() string {
	return "commands-checker"
}

func (c *CommandsChecker) Run(ctx context.Context) (context.Context, error) {
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get commands from context failed")
	}

	if len(commands) == 0 {
		return ctx, errors.New("no commands found")
	}

	return ctx, nil
}

func NewCommandsChecker() *CommandsChecker {
	return &CommandsChecker{}
}
