package core

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
)

type BinaryRemover struct {
	BaseStep
}

func (s *BinaryRemover) String() string {
	return "binary-remover"
}

func (s *BinaryRemover) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	fs := define.FS

	command, err := GetCommandFromContext(ctx)
	if err != nil {
		return ctx, nil
	}

	if !command.Managed {
		return ctx, nil
	}

	logger.Info("removing binary", map[string]interface{}{
		"location": command.Location,
	})

	err = fs.Remove(command.Location)
	if err != nil {
		return ctx, errors.Wrapf(err, "remove binary failed")
	}

	return ctx, nil
}

func NewBinaryRemover() *BinaryRemover {
	return &BinaryRemover{}
}

type CommandRemover struct {
	BaseStep
}

func (s *CommandRemover) String() string {
	return "command-remover"
}

func (s *CommandRemover) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	client := GetDBClientFromContext(ctx)

	command, err := GetCommandFromContext(ctx)
	if err != nil {
		return ctx, nil
	}

	logger.Info("removing command", map[string]interface{}{
		"name":    command.Name,
		"version": command.Version,
	})

	err = client.DeleteStruct(command)
	if err != nil {
		return ctx, errors.Wrapf(err, "remove command failed")
	}

	return ctx, nil
}

func NewCommandRemover() *CommandRemover {
	return &CommandRemover{}
}

type BrokenCommandsRemover struct {
	BaseStep
}

func (s *BrokenCommandsRemover) String() string {
	return "broken-commands-remover"
}

func (s *BrokenCommandsRemover) Run(ctx context.Context) (context.Context, error) {
	fs := define.FS
	logger := define.Logger
	var errs error
	client := GetDBClientFromContext(ctx)
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	availableCommands := make([]*model.Command, 0, len(commands))
	for _, command := range commands {
		location := command.Location
		if command.Managed {
			location = GetCommandPath(command.Name, command.Version)
		}

		_, err := fs.Stat(location)
		if err == nil {
			availableCommands = append(availableCommands, command)
			continue
		}

		logger.Debug("deleting command", map[string]interface{}{
			"name":     command.Name,
			"version":  command.Version,
			"location": command.Location,
			"err":      err,
		})

		err = client.DeleteStruct(command)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "remove command %s(%s) failed", command.Name, command.Version))
		}
		logger.Info("command deleted", map[string]interface{}{
			"name":    command.Name,
			"version": command.Version,
		})

	}

	return ctx, errs
}

func NewBrokenCommandsRemover() *BrokenCommandsRemover {
	return &BrokenCommandsRemover{}
}
