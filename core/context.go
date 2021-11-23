package core

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

func GetDBClientFromContext(ctx context.Context) DBClient {
	return ctx.Value(define.ContextKeyDBClient).(DBClient)
}

func getCommandFromContext(ctx context.Context) *model.Command {
	value := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommand)
	if value == nil {
		return nil
	}

	command, ok := value.(*model.Command)
	if !ok {
		return nil
	}

	return command
}

func getCommandsFromContext(ctx context.Context) []*model.Command {
	values := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommands)
	if values == nil {
		return nil
	}

	commands, ok := values.([]*model.Command)
	if !ok || len(commands) == 0 {
		return nil
	}

	return commands
}

func GetCommandFromContext(ctx context.Context) (*model.Command, error) {
	command := getCommandFromContext(ctx)
	if command != nil {
		return command, nil
	}

	commands := getCommandsFromContext(ctx)
	if commands != nil && len(commands) > 0 {
		return commands[0], nil
	}

	return nil, errors.Wrapf(ErrContextValueNotFound, "command not found")
}

func GetCommandsFromContext(ctx context.Context) ([]*model.Command, error) {
	command := getCommandFromContext(ctx)
	if command != nil {
		return []*model.Command{command}, nil
	}

	commands := getCommandsFromContext(ctx)
	if commands != nil && len(commands) > 0 {
		return commands, nil
	}

	return nil, errors.Wrapf(ErrContextValueNotFound, "commands not found")
}

type ContextValueSetter struct {
	BaseStep
	values map[define.ContextKey]interface{}
}

func (s *ContextValueSetter) String() string {
	return "context-value-setter"
}

func (s *ContextValueSetter) Run(ctx context.Context) (context.Context, error) {
	return utils.SetIntoContext(ctx, s.values), nil
}

func NewContextValueSetter(values map[define.ContextKey]interface{}) *ContextValueSetter {
	return &ContextValueSetter{
		values: values,
	}
}
