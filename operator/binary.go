package operator

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
)

type BinariesInstaller struct {
	*CmdrOperator
	managed bool
}

func (i *BinariesInstaller) String() string {
	return "binaries-installer"
}

func (i *BinariesInstaller) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get commands from context failed")
	}

	var errs error
	for _, command := range commands {
		name := command.Name
		version := command.Version
		location := command.Location

		logger.Info("installing binary", map[string]interface{}{
			"name":     name,
			"version":  version,
			"location": location,
		})

		err = i.cmdr.BinaryManager.Install(name, version, location, !i.managed)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "install %s(%s) binary failed", name, version))
			continue
		}
	}

	return ctx, errors.Wrap(errs, "install binaries failed")
}

func NewBinariesInstaller(cmdr *core.Cmdr, managed bool) *BinariesInstaller {
	return &BinariesInstaller{
		CmdrOperator: NewCmdrOperator(cmdr),
		managed:      managed,
	}
}

type BinariesUninstaller struct {
	*CmdrOperator
}

func (i *BinariesUninstaller) String() string {
	return "binaries-uninstaller"
}

func (i *BinariesUninstaller) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	var errs error
	for _, command := range commands {
		logger.Info("uninstalling binary", map[string]interface{}{
			"location": command.Location,
		})

		err = i.cmdr.BinaryManager.Uninstall(command.Name, command.Version)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, command.Location))
		}
	}

	return ctx, errs
}

func NewBinariesUninstaller(cmdr *core.Cmdr) *BinariesUninstaller {
	return &BinariesUninstaller{
		CmdrOperator: NewCmdrOperator(cmdr),
	}
}

type BinariesActivator struct {
	*CmdrOperator
}

func (s *BinariesActivator) String() string {
	return "binaries-activator"
}

func (s *BinariesActivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "get commands from context failed")
	}

	logger.Info("activating binaries", map[string]interface{}{
		"count": len(commands),
	})

	var errs error
	for _, command := range commands {
		err = s.cmdr.BinaryManager.Activate(command.Name, command.Version)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "activate %s(%s) binary failed", command.Name, command.Version))
			continue
		}
	}

	return ctx, errs
}

func NewBinariesActivator(cmdr *core.Cmdr) *BinariesActivator {
	return &BinariesActivator{
		CmdrOperator: NewCmdrOperator(cmdr),
	}
}

type BinariesDeactivator struct {
	*CmdrOperator
}

func (s *BinariesDeactivator) String() string {
	return "binaries-deactivator"
}

func (s *BinariesDeactivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "get commands from context failed")
	}

	logger.Info("deactivating binaries", map[string]interface{}{
		"count": len(commands),
	})

	var errs error
	for _, command := range commands {
		err = s.cmdr.BinaryManager.Deactivate(command.Name)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "remove %s failed", command.Name))
			continue
		}
	}

	return ctx, errs
}

func NewBinariesDeactivator(cmdr *core.Cmdr) *BinariesDeactivator {
	return &BinariesDeactivator{
		CmdrOperator: NewCmdrOperator(cmdr),
	}
}
