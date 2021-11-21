package core

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

func installBinary(name, version, location string) error {
	fs := define.FS
	logger := define.Logger
	dir := GetCommandDir(name)
	target := GetCommandPath(name, version)

	logger.Info("installing binary", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})

	logger.Debug("creating binary dir", map[string]interface{}{
		"name":     name,
		"location": location,
		"dir":      dir,
	})
	err := fs.MkdirAll(dir, 0755)
	if err != nil {
		return errors.Wrapf(err, "create dir %s failed", dir)
	}

	logger.Debug("coping command", map[string]interface{}{
		"name":     name,
		"location": location,
		"target":   target,
	})
	err = utils.CopyFile(location, target)
	if err != nil {
		return errors.WithMessagef(err, "install command %s failed", target)
	}

	err = fs.Chmod(target, 0755)
	if err != nil {
		return errors.WithMessagef(err, "change command mode %s failed", target)
	}

	return nil
}

type BinaryInstaller struct {
	BaseStep
}

func (i *BinaryInstaller) String() string {
	return "binary-installer"
}

func (i *BinaryInstaller) Run(ctx context.Context) (context.Context, error) {
	name := utils.GetStringFromContext(ctx, define.ContextKeyName)
	version := utils.GetStringFromContext(ctx, define.ContextKeyVersion)
	location := utils.GetStringFromContext(ctx, define.ContextKeyLocation)

	err := installBinary(name, version, location)
	if err != nil {
		return ctx, errors.WithMessagef(err, "install command %s failed", name)
	}

	return ctx, nil
}

func NewBinaryInstaller() *BinaryInstaller {
	return &BinaryInstaller{}
}

type BinariesInstaller struct {
	BaseStep
}

func (i *BinariesInstaller) String() string {
	return "binaries-installer"
}

func (i *BinariesInstaller) Run(ctx context.Context) (context.Context, error) {
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, nil
	}

	var errs error
	for _, command := range commands {
		err := installBinary(command.Name, command.Version, command.Location)
		if err != nil {
			errs = multierror.Append(errs, err)
		}

	}

	return ctx, errs
}

func NewBinariesInstaller() *BinariesInstaller {
	return &BinariesInstaller{}
}

type CommandInstaller struct {
	BaseStep
}

func (i *CommandInstaller) String() string {
	return "command-installer"
}

func (i *CommandInstaller) Run(ctx context.Context) (context.Context, error) {
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

	logger.Info("installing command", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
		"managed":  managed,
	})

	err := client.Save(&model.Command{
		Name:     name,
		Version:  version,
		Location: location,
		Managed:  managed,
	})

	if err != nil {
		return ctx, errors.Wrapf(err, "create command failed")
	}

	return ctx, nil
}

func NewCommandInstaller() *CommandInstaller {
	return &CommandInstaller{}
}
