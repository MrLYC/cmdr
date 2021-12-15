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
	commands := getCommandsFromContext(ctx)
	if commands == nil {
		return nil, errors.Wrapf(ErrContextValueNotFound, "commands not found")
	}
	return commands[0], nil
}

func GetCommandsFromContext(ctx context.Context) ([]*model.Command, error) {
	commands := getCommandsFromContext(ctx)
	if commands == nil {
		return nil, errors.Wrapf(ErrContextValueNotFound, "commands not found")
	}
	return commands, nil

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
