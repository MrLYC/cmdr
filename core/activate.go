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

type CommandActivator struct {
	BaseStep
}

func (s *CommandActivator) String() string {
	return "command-activator"
}

func (s *CommandActivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	command, err := GetCommandFromContext(ctx)
	if err != nil {
		return ctx, nil
	}

	logger.Info("activating command", map[string]interface{}{
		"name":    command.Name,
		"version": command.Version,
	})

	client := GetDBClientFromContext(ctx)
	command.Activated = true
	err = client.Save(command)
	if err != nil {
		return ctx, errors.Wrapf(err, "save command failed")
	}

	return ctx, nil
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
	var commands []*model.Command

	err := client.Select(q.Eq("Name", name), q.Eq("Activated", true)).Find(&commands)
	if errors.Cause(err) == storm.ErrNotFound {
		return ctx, nil
	} else if err != nil {
		return ctx, errors.Wrapf(err, "get command failed")
	}

	for _, command := range commands {
		logger.Info("deactivating command", map[string]interface{}{
			"name":    command.Name,
			"version": command.Version,
		})

		command.Activated = false
		err = client.Save(command)
		if err != nil {
			return ctx, errors.Wrapf(err, "save command failed")
		}
	}

	return ctx, nil
}

func NewCommandDeactivator() *CommandDeactivator {
	return &CommandDeactivator{}
}
