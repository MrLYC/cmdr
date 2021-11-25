package core

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

type CommandDefiner struct {
	BaseStep
}

func (i *CommandDefiner) String() string {
	return "command-definer"
}

func (i *CommandDefiner) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	name := utils.GetStringFromContext(ctx, define.ContextKeyName)
	version := utils.GetStringFromContext(ctx, define.ContextKeyVersion)
	managed := utils.GetBoolFromContext(ctx, define.ContextKeyCommandManaged)
	client := GetDBClientFromContext(ctx)
	var location string
	if managed {
		location = GetCommandPath(name, version)
	} else {
		location = utils.GetStringFromContext(ctx, define.ContextKeyLocation)
	}

	logger.Info("define command", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
		"managed":  managed,
	})

	command := model.Command{
		Name:     name,
		Version:  version,
		Location: location,
		Managed:  managed,
	}

	err := client.Save(&command)

	if err != nil {
		return ctx, errors.Wrapf(err, "create command failed")
	}

	return context.WithValue(ctx, define.ContextKeyCommands, []*model.Command{&command}), nil
}

func NewCommandDefiner() *CommandDefiner {
	return &CommandDefiner{}
}

type CommandUndefiner struct {
	BaseStep
}

func (s *CommandUndefiner) String() string {
	return "command-undefiner"
}

func (s *CommandUndefiner) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	client := GetDBClientFromContext(ctx)

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, nil
	}

	for _, command := range commands {
		logger.Info("undefine command", map[string]interface{}{
			"name":    command.Name,
			"version": command.Version,
		})

		err = client.DeleteStruct(command)
		if err != nil {
			return ctx, errors.Wrapf(err, "remove command failed")
		}
	}

	return ctx, nil
}

func NewCommandUndefiner() *CommandUndefiner {
	return &CommandUndefiner{}
}
