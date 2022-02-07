package operator

import (
	"context"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type BinariesInstaller struct {
	BaseOperator
	helper *utils.CmdrHelper
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
		if !command.Managed {
			continue
		}

		name := command.Name
		version := command.Version
		location := command.Location
		dir := i.helper.GetCommandShimsDir(name)
		target := i.helper.GetCommandShimsPath(name, version)

		logger.Info("installing binary", map[string]interface{}{
			"name":     name,
			"version":  version,
			"location": location,
		})

		logger.Debug("creating binary dir", map[string]interface{}{
			"dir": dir,
		})
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "create dir %s failed", dir))
			continue
		}

		logger.Debug("coping command", map[string]interface{}{
			"name":     name,
			"location": location,
			"target":   target,
		})
		err = utils.CopyFile(location, target)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "install command %s failed", target))
			continue
		}

		err = os.Chmod(target, 0755)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "change command mode %s failed", target))
			continue
		}
	}

	return ctx, errors.Wrap(errs, "install binaries failed")
}

func NewBinariesInstaller(helper *utils.CmdrHelper) *BinariesInstaller {
	return &BinariesInstaller{
		helper: helper,
	}
}

type BinariesUninstaller struct {
	BaseOperator
}

func (s *BinariesUninstaller) String() string {
	return "binaries-uninstaller"
}

func (s *BinariesUninstaller) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	var errs error
	for _, command := range commands {
		if !command.Managed {
			continue
		}

		_, err := os.Stat(command.Location)
		if err != nil {
			logger.Debug("binary not found", map[string]interface{}{
				"location": command.Location,
			})
			continue
		}

		logger.Info("removing binary", map[string]interface{}{
			"location": command.Location,
		})

		err = os.Remove(command.Location)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, command.Location))
		}
	}

	return ctx, errs
}

func NewBinariesUninstaller() *BinariesUninstaller {
	return &BinariesUninstaller{}
}

type BinariesActivator struct {
	BaseOperator
	helper *utils.CmdrHelper
}

func (s *BinariesActivator) String() string {
	return "binaries-activator"
}

func (s *BinariesActivator) cleanUpBinary(binPath string) {
	logger := define.Logger

	info, err := os.Lstat(binPath)
	if err != nil {
		return
	}

	logger.Debug("remove exists binary", map[string]interface{}{
		"path": binPath,
	})

	if info.IsDir() {
		_ = os.RemoveAll(binPath)
	} else {
		_ = os.Remove(binPath)
	}
}

func (s *BinariesActivator) activateBinary(name, location string) error {
	binPath := s.helper.GetCommandBinPath(name)
	s.cleanUpBinary(binPath)

	err := os.Symlink(location, binPath)
	if err != nil {
		return errors.Wrapf(err, "create symbol link failed")
	}

	return nil
}

func (s *BinariesActivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	binDir := s.helper.GetBinDir()

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "get commands from context failed")
	}

	logger.Info("activating binaries", map[string]interface{}{
		"count": len(commands),
	})

	err = os.MkdirAll(binDir, 0755)
	if err != nil {
		return ctx, errors.Wrapf(err, "create bin dir %s failed", binDir)
	}

	var errs error
	for _, command := range commands {
		location := command.Location
		if command.Managed {
			location = s.helper.GetCommandShimsPath(command.Name, command.Version)
		}

		err = s.activateBinary(command.Name, location)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "activate %s(%s) binary failed", command.Name, command.Version))
			continue
		}
	}

	return ctx, errs
}

func NewBinariesActivator(helper *utils.CmdrHelper) *BinariesActivator {
	return &BinariesActivator{
		helper: helper,
	}
}

type BinariesDeactivator struct {
	BaseOperator
	helper *utils.CmdrHelper
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
		binPath := s.helper.GetCommandBinPath(command.Name)
		_, err := os.Stat(binPath)
		if err != nil {
			logger.Debug("binary not found", map[string]interface{}{
				"path": binPath,
			})
			continue
		}

		logger.Info("deactivating binary", map[string]interface{}{
			"path": binPath,
		})

		err = os.Remove(binPath)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "remove %s failed", binPath))
			continue
		}
	}

	return ctx, errs
}

func NewBinariesDeactivator(helper *utils.CmdrHelper) *BinariesDeactivator {
	return &BinariesDeactivator{
		helper: helper,
	}
}
