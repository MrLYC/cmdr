package core

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
)

type CommandDefiner struct {
	BaseStep
	shimsDir string
	command  model.Command
}

func (i *CommandDefiner) String() string {
	return "command-definer"
}

func (i *CommandDefiner) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	client := GetDBClientFromContext(ctx)

	logger.Info("define command", map[string]interface{}{
		"name":     i.command.Name,
		"version":  i.command.Version,
		"location": i.command.Location,
		"managed":  i.command.Managed,
	})

	var command model.Command
	err := client.Select(q.Eq("Name", i.command.Name), q.Eq("Version", i.command.Version)).First(&command)
	switch errors.Cause(err) {
	case nil:
		i.command.ID = command.ID
	case storm.ErrNotFound:
	default:
		return ctx, errors.Wrapf(err, "define command failed")
	}

	return context.WithValue(ctx, define.ContextKeyCommands, []*model.Command{&i.command}), nil
}

func (i *CommandDefiner) Commit(ctx context.Context) error {
	logger := define.Logger
	client := GetDBClientFromContext(ctx)

	if i.command.Managed {
		i.command.Location = GetCommandPath(i.shimsDir, i.command.Name, i.command.Version)
	}

	logger.Debug("saving command", map[string]interface{}{
		"name":     i.command.Name,
		"version":  i.command.Version,
		"location": i.command.Location,
		"managed":  i.command.Managed,
	})

	err := client.Save(&i.command)
	if err != nil {
		return errors.Wrapf(err, "save command failed")
	}

	return nil
}

func NewCommandDefiner(shimsDir, name, version, location string, managed bool) *CommandDefiner {
	return &CommandDefiner{
		shimsDir: shimsDir,
		command: model.Command{
			Name:     name,
			Version:  version,
			Location: location,
			Managed:  managed,
		},
	}
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
		return ctx, err
	}

	var errs error
	for _, command := range commands {
		logger.Info("undefine command", map[string]interface{}{
			"name":    command.Name,
			"version": command.Version,
		})

		err = client.DeleteStruct(command)
		switch err {
		case nil, storm.ErrNotFound:
		default:
			errs = multierror.Append(errs, errors.Wrapf(err, "delete command failed"))
		}
	}

	return ctx, errs
}

func NewCommandUndefiner() *CommandUndefiner {
	return &CommandUndefiner{}
}

type CommandActivator struct {
	BaseStep
}

func (s *CommandActivator) String() string {
	return "command-activator"
}

func (s *CommandActivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	var errs error
	for _, command := range commands {
		logger.Info("activating command", map[string]interface{}{
			"name":    command.Name,
			"version": command.Version,
		})

		client := GetDBClientFromContext(ctx)
		command.Activated = true
		err = client.Save(command)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "activate command failed"))
		}
	}

	return ctx, errs
}

func NewCommandActivator() *CommandActivator {
	return &CommandActivator{}
}

type CommandsDeactivator struct {
	BaseStep
}

func (s *CommandsDeactivator) String() string {
	return "command-deactivator"
}

func (s *CommandsDeactivator) deactivateCommand(ctx context.Context, command *model.Command) error {
	logger := define.Logger
	client := GetDBClientFromContext(ctx)
	var queriedCommands []*model.Command

	err := client.Select(q.Eq("Name", command.Name), q.Eq("Activated", true)).Find(&queriedCommands)
	switch err {
	case nil:
	case storm.ErrNotFound:
		return nil
	default:
		return errors.Wrapf(err, "query command %s failed", command.Name)
	}

	var errs error
	for _, queried := range queriedCommands {
		if queried.Name == command.Name && queried.Version == command.Version {
			continue
		}

		logger.Info("deactivating command", map[string]interface{}{
			"name":    queried.Name,
			"version": queried.Version,
		})

		queried.Activated = false
		err = client.Save(queried)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "deactivate command failed"))
		}
	}

	return errs
}

func (s *CommandsDeactivator) Run(ctx context.Context) (context.Context, error) {
	var errs error

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get commands from context failed")
	}

	for _, command := range commands {
		err := s.deactivateCommand(ctx, command)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return ctx, errs
}

func NewCommandDeactivator() *CommandsDeactivator {
	return &CommandsDeactivator{}
}
