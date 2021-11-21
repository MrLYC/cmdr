package core

import (
	"context"
	"path/filepath"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
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
	value := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommand)
	if value == nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "command not found")
	}

	command, ok := value.(*model.Command)
	if !ok || command == nil {
		return ctx, nil
	}

	logger.Info("activating command", map[string]interface{}{
		"name":    command.Name,
		"version": command.Version,
	})

	client := GetDBClientFromContext(ctx)
	command.Activated = true
	err := client.Save(command)
	if err != nil {
		return ctx, errors.Wrapf(err, "save command failed")
	}

	return ctx, nil
}

func NewCommandActivator() *CommandActivator {
	return &CommandActivator{}
}

type BinaryActivator struct {
	BaseStep
}

func (s *BinaryActivator) String() string {
	return "binary-activator"
}

func (s *BinaryActivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	value := utils.GetInterfaceFromContext(ctx, define.ContextKeyCommand)
	if value == nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "command not found")
	}

	command, ok := value.(*model.Command)
	if !ok || command == nil {
		return ctx, nil
	}

	logger.Info("activating binary", map[string]interface{}{
		"name":    command.Name,
		"version": command.Version,
	})

	fs := define.FS
	binPath := filepath.Join(GetBinDir(), command.Name)

	linkReader := define.GetSymbolLinkReader()
	_, err := linkReader.ReadlinkIfPossible(binPath)
	if err == nil {
		logger.Debug("remove exists binary", map[string]interface{}{
			"name":   command.Name,
			"target": command.Location,
		})
		fs.Remove(binPath)
	}

	linker := define.GetSymbolLinker()
	err = linker.SymlinkIfPossible(command.Location, binPath)
	if err != nil {
		return ctx, errors.Wrapf(err, "create symbol link failed")
	}

	return ctx, nil
}

func NewBinaryActivator() *BinaryActivator {
	return &BinaryActivator{}
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
