package core

import (
	"context"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type BinariesInstaller struct {
	BaseStep
	shimsDir string
}

func (i *BinariesInstaller) String() string {
	return "binaries-installer"
}

func (i *BinariesInstaller) Run(ctx context.Context) (context.Context, error) {
	fs := define.FS
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
		dir := GetCommandDir(i.shimsDir, name)
		target := GetCommandPath(i.shimsDir, name, version)

		logger.Info("installing binary", map[string]interface{}{
			"name":     name,
			"version":  version,
			"location": location,
		})

		logger.Debug("creating binary dir", map[string]interface{}{
			"dir": dir,
		})
		err := fs.MkdirAll(dir, 0755)
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

		err = fs.Chmod(target, 0755)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "change command mode %s failed", target))
			continue
		}
	}

	return ctx, errors.Wrap(errs, "install binaries failed")
}

func NewBinariesInstaller(shimsDir string) *BinariesInstaller {
	return &BinariesInstaller{
		shimsDir: shimsDir,
	}
}

type BinariesUninstaller struct {
	BaseStep
}

func (s *BinariesUninstaller) String() string {
	return "binaries-uninstaller"
}

func (s *BinariesUninstaller) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	fs := define.FS

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	var errs error
	for _, command := range commands {
		if !command.Managed {
			continue
		}

		exists, err := afero.Exists(fs, command.Location)
		if !exists || err != nil {
			logger.Debug("binary not found", map[string]interface{}{
				"location": command.Location,
			})
			continue
		}

		logger.Info("removing binary", map[string]interface{}{
			"location": command.Location,
		})

		err = fs.Remove(command.Location)
		if err != nil {
			multierror.Append(errs, errors.WithMessagef(err, command.Location))
			continue
		}
	}

	return ctx, errs
}

func NewBinariesUninstaller() *BinariesUninstaller {
	return &BinariesUninstaller{}
}

type BinariesActivator struct {
	BaseStep
	binDir string
}

func (s *BinariesActivator) String() string {
	return "binaries-activator"
}

func (s *BinariesActivator) cleanUpBinary(binPath string) {
	fs := define.FS
	logger := define.Logger

	info, err := fs.Stat(binPath)
	if err != nil {
		return
	}

	logger.Debug("remove exists binary", map[string]interface{}{
		"path": binPath,
	})

	if info.IsDir() {
		_ = fs.RemoveAll(binPath)
	} else {
		_ = fs.Remove(binPath)
	}
}

func (s *BinariesActivator) activateBinary(name, location string) error {
	binPath := filepath.Join(s.binDir, name)
	s.cleanUpBinary(binPath)

	linker := define.GetSymbolLinker()
	err := linker.SymlinkIfPossible(location, binPath)
	if err != nil {
		return errors.Wrapf(err, "create symbol link failed")
	}

	return nil
}

func (s *BinariesActivator) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	fs := define.FS

	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(ErrContextValueNotFound, "get commands from context failed")
	}

	logger.Info("activating binaries", map[string]interface{}{
		"count": len(commands),
	})

	err = fs.MkdirAll(s.binDir, 0755)
	if err != nil {
		return ctx, errors.Wrapf(err, "create bin dir %s failed", s.binDir)
	}

	var errs error
	for _, command := range commands {
		err = s.activateBinary(command.Name, command.Location)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "activate %s(%s) binary failed", command.Name, command.Version))
			continue
		}
	}

	return ctx, errs
}

func NewBinariesActivator(binDir string) *BinariesActivator {
	return &BinariesActivator{
		binDir: binDir,
	}
}
