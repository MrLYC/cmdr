package core

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
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

	value := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommand)
	if value == nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "command not found")
	}

	command, ok := value.(*model.Command)
	if !ok || command == nil {
		return ctx, nil
	}

	if !command.Managed {
		return ctx, nil
	}

	logger.Info("removing binary", map[string]interface{}{
		"location": command.Location,
	})

	err := fs.Remove(command.Location)
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

	value := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommand)
	if value == nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "command not found")
	}

	command, ok := value.(*model.Command)
	if !ok || command == nil {
		return ctx, nil
	}

	logger.Info("removing command", map[string]interface{}{
		"name":    command.Name,
		"version": command.Version,
	})

	err := client.DeleteStruct(command)
	if err != nil {
		return ctx, errors.Wrapf(err, "remove command failed")
	}

	return ctx, nil
}

func NewCommandRemover() *CommandRemover {
	return &CommandRemover{}
}
