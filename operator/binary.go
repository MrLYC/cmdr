package operator

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type BinariesInstaller struct {
	BaseOperator
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
		dir := utils.GetCommandShimsDir(i.shimsDir, name)
		target := utils.GetCommandShimsPath(i.shimsDir, name, version)

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
	BaseOperator
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
	BaseOperator
	shimsDir string
	binDir   string
}

func (s *BinariesActivator) String() string {
	return "binaries-activator"
}

func (s *BinariesActivator) cleanUpBinary(binPath string) {
	fs := define.FS
	lstater := utils.GetFsLstater()
	logger := define.Logger

	info, _, err := lstater.LstatIfPossible(binPath)
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
	binPath := utils.GetCommandBinPath(s.binDir, name)
	s.cleanUpBinary(binPath)

	linker := utils.GetSymbolLinker()
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
		location := command.Location
		if command.Managed {
			location = utils.GetCommandShimsPath(s.shimsDir, command.Name, command.Version)
		}

		err = s.activateBinary(command.Name, location)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "activate %s(%s) binary failed", command.Name, command.Version))
			continue
		}
	}

	return ctx, errs
}

func NewBinariesActivator(binDir, shimsDir string) *BinariesActivator {
	return &BinariesActivator{
		binDir:   binDir,
		shimsDir: shimsDir,
	}
}
