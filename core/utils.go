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

func GetCommandFromContext(ctx context.Context) (*model.Command, error) {
	values := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommands)
	if values == nil {
		return nil, errors.Wrapf(ErrContextValueNotFound, "commands not found")
	}

	command, ok := values.(*model.Command)
	if !ok {
		return nil, nil
	}

	return command, nil
}

func GetCommandsFromContext(ctx context.Context) ([]*model.Command, error) {
	values := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommands)
	if values == nil {
		return nil, errors.Wrapf(ErrContextValueNotFound, "commands not found")
	}

	commands, ok := values.([]*model.Command)
	if !ok || len(commands) == 0 {
		return nil, nil
	}

	return commands, nil
}
