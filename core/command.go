package core

import (
	"context"
	"path/filepath"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

type CommandDefiner struct {
	BaseStep
	shimsDir string
}

func (i *CommandDefiner) String() string {
	return "command-definer"
}

func (i *CommandDefiner) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	name := utils.GetStringFromContext(ctx, define.ContextKeyName)
	version := utils.GetStringFromContext(ctx, define.ContextKeyVersion)
	if name == "" {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "name is empty")
	} else if version == "" {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "version is empty")
	}

	managed := utils.GetBoolFromContext(ctx, define.ContextKeyCommandManaged)
	client := GetDBClientFromContext(ctx)
	var location string
	if managed {
		location = GetCommandPath(i.shimsDir, name, version)
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
		Name:    name,
		Version: version,
	}

	err := client.Select(q.Eq("name", name), q.Eq("version", version)).First(&command)
	switch errors.Cause(err) {
	case nil, storm.ErrNotFound:
		command.Location = location
		command.Managed = managed
	default:
		return ctx, errors.Wrapf(err, "define command failed")
	}

	err = client.Save(&command)
	if err != nil {
		return ctx, errors.Wrapf(err, "update command failed")
	}

	return context.WithValue(ctx, define.ContextKeyCommands, []*model.Command{&command}), nil
}

func NewCommandDefiner(shimsDir string) *CommandDefiner {
	return &CommandDefiner{
		shimsDir: shimsDir,
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

func activateBinary(name, location string) error {
	fs := define.FS
	logger := define.Logger
	binPath := filepath.Join(GetBinDir(), name)

	linkReader := define.GetSymbolLinkReader()
	_, err := linkReader.ReadlinkIfPossible(binPath)
	if err == nil {
		logger.Debug("remove exists binary", map[string]interface{}{
			"name":   name,
			"target": location,
		})
		fs.Remove(binPath)
	}

	linker := define.GetSymbolLinker()
	err = linker.SymlinkIfPossible(location, binPath)
	if err != nil {
		return errors.Wrapf(err, "create symbol link failed")
	}

	return nil
}

type BinaryActivator struct {
	BaseStep
}

func (s *BinaryActivator) String() string {
	return "binary-activator"
}

func (s *BinaryActivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	command, err := GetCommandFromContext(ctx)
	if err != nil {
		return ctx, nil
	}

	logger.Info("activating binary", map[string]interface{}{
		"name":    command.Name,
		"version": command.Version,
	})

	err = activateBinary(command.Name, command.Location)
	if err != nil {
		return ctx, errors.Wrapf(err, "activate binary failed")
	}

	return ctx, nil
}

func NewBinaryActivator() *BinaryActivator {
	return &BinaryActivator{}
}

type BinariesActivator struct {
	BaseStep
}

func (s *BinariesActivator) String() string {
	return "binaries-activator"
}

func (s *BinariesActivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, nil
	}

	logger.Info("activating binaries", map[string]interface{}{
		"count": len(commands),
	})

	var errs error
	for _, command := range commands {
		err = activateBinary(command.Name, command.Location)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "activate %s(%s) binary failed", command.Name, command.Version))
			continue
		}
	}

	return ctx, nil
}

func NewBinariesActivator() *BinariesActivator {
	return &BinariesActivator{}
}

type CommandDeactivator struct {
	BaseStep
}

func (s *CommandDeactivator) String() string {
	return "command-deactivator"
}

func (s *CommandDeactivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	client := GetDBClientFromContext(ctx)
	name := utils.GetStringFromContext(ctx, define.ContextKeyName)
	if name == "" {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "name is empty")
	}

	var commands []*model.Command

	err := client.Select(q.Eq("Name", name), q.Eq("Activated", true)).Find(&commands)
	switch err {
	case nil:
	case storm.ErrNotFound:
		return ctx, nil
	default:
		return ctx, errors.Wrapf(err, "get command failed")
	}

	var errs error
	for _, command := range commands {
		logger.Info("deactivating command", map[string]interface{}{
			"name":    command.Name,
			"version": command.Version,
		})

		command.Activated = false
		err = client.Save(command)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "deactivate command failed"))
		}
	}

	return ctx, errs
}

func NewCommandDeactivator() *CommandDeactivator {
	return &CommandDeactivator{}
}
