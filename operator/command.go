package operator

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/model"
	"github.com/mrlyc/cmdr/define"
)

type CommandDefiner struct {
	*CmdrOperator
	command model.Command
}

func (i *CommandDefiner) String() string {
	return "command-definer"
}

func (i *CommandDefiner) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	logger.Info("define command", map[string]interface{}{
		"name":     i.command.Name,
		"version":  i.command.Version,
		"location": i.command.Location,
	})

	var command model.Command
	err := i.cmdr.CommandManager.
		Query().
		WithName(i.command.Name).
		WithVersion(i.command.Version).
		Done().
		First(&command)

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

	logger.Debug("saving command", map[string]interface{}{
		"name":     i.command.Name,
		"version":  i.command.Version,
		"location": i.command.Location,
	})

	err := i.cmdr.CommandManager.Client.Update(&i.command)
	if err != nil {
		return errors.Wrapf(err, "save command failed")
	}

	return nil
}

func NewCommandDefiner(cmdr *core.Cmdr, name, version, location string) *CommandDefiner {
	return &CommandDefiner{
		CmdrOperator: NewCmdrOperator(cmdr),
		command: model.Command{
			Name:     name,
			Version:  version,
			Location: location,
		},
	}
}

type CommandUndefiner struct {
	*CmdrOperator
}

func (i *CommandUndefiner) String() string {
	return "command-undefiner"
}

func (i *CommandUndefiner) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

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

		err = i.cmdr.CommandManager.Client.DeleteStruct(command)
		switch err {
		case nil, storm.ErrNotFound:
		default:
			errs = multierror.Append(errs, errors.Wrapf(err, "delete command failed"))
		}
	}

	return ctx, errs
}

func NewCommandUndefiner(cmdr *core.Cmdr) *CommandUndefiner {
	return &CommandUndefiner{
		CmdrOperator: NewCmdrOperator(cmdr),
	}
}

type CommandActivator struct {
	*CmdrOperator
}

func (i *CommandActivator) String() string {
	return "command-activator"
}

func (i *CommandActivator) Run(ctx context.Context) (context.Context, error) {
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

		err = i.cmdr.CommandManager.Activate(command)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "activate command failed"))
		}
	}

	return ctx, errs
}

func NewCommandActivator(cmdr *core.Cmdr) *CommandActivator {
	return &CommandActivator{
		CmdrOperator: NewCmdrOperator(cmdr),
	}
}

type CommandsDeactivator struct {
	*CmdrOperator
}

func (i *CommandsDeactivator) String() string {
	return "command-deactivator"
}

func (i *CommandsDeactivator) Run(ctx context.Context) (context.Context, error) {
	var errs error

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get commands from context failed")
	}

	for _, command := range commands {
		err := i.cmdr.CommandManager.DeactivateAll(command.Name)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return ctx, errs
}

func NewCommandDeactivator(cmdr *core.Cmdr) *CommandsDeactivator {
	return &CommandsDeactivator{
		CmdrOperator: NewCmdrOperator(cmdr),
	}
}

type CommandHandler struct {
	BaseOperator
	name   string
	runner func(context.Context, []*model.Command) error
}

func (c *CommandHandler) String() string {
	return c.name
}

func (c *CommandHandler) Run(ctx context.Context) (context.Context, error) {
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get commands from context failed")
	}

	return ctx, c.runner(ctx, commands)
}

func NewCommandHandler(name string, runner func(ctx context.Context, commands []*model.Command) error) *CommandHandler {
	return &CommandHandler{
		name:   name,
		runner: runner,
	}
}
