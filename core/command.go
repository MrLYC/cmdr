package core

import (
	"context"
	"fmt"
	"path"

	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/utils"
)

type CommandHelper struct {
	client   Client
	shimsDir string
	binDir   string
}

func (h *CommandHelper) GetCommandByNameAndVersion(ctx context.Context, name, version string) (*model.Command, error) {
	return h.GetCommand(ctx, q.Eq("Name", name), q.Eq("Version", version))
}

func (h *CommandHelper) defineCommand(ctx context.Context, name, version, location string, managed bool) error {
	logger := define.Logger

	logger.Debug("saving command record", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})

	err := h.client.Save(&model.Command{
		Name:     name,
		Version:  version,
		Location: location,
		Managed:  managed,
	})

	if err != nil {
		return errors.Wrapf(err, "create command failed")
	}

	return nil
}

func (h *CommandHelper) Define(ctx context.Context, name, version, location string) error {
	return h.client.Atomic(func() error {
		logger := define.Logger

		logger.Info("checking command", map[string]interface{}{
			"name":     name,
			"version":  version,
			"location": location,
		})

		command, err := h.GetCommandByNameAndVersion(ctx, name, version)
		if err == nil && command != nil {
			logger.Debug("command exists", map[string]interface{}{
				"name":     name,
				"version":  version,
				"location": command.Location,
			})

			return errors.Wrapf(ErrCommandAlreadyExists, name)
		}

		return h.defineCommand(ctx, name, version, location, false)
	})
}

func (h *CommandHelper) installCommandBinary(name, version, location, target string) error {
	fs := define.FS
	logger := define.Logger

	dir := path.Dir(target)
	logger.Debug("creating binary dir", map[string]interface{}{
		"name":     name,
		"location": location,
		"dir":      dir,
	})
	utils.CheckError(fs.MkdirAll(dir, 0755))

	logger.Debug("coping command", map[string]interface{}{
		"name":     name,
		"location": location,
		"target":   target,
	})
	err := utils.CopyFile(location, target)
	if err != nil {
		return errors.WithMessagef(err, "install command %s failed", target)
	}

	err = fs.Chmod(target, 0755)
	if err != nil {
		return errors.WithMessagef(err, "change command mode %s failed", target)
	}

	return nil
}

func (h *CommandHelper) Install(ctx context.Context, name, version, location string) error {
	return h.client.Atomic(func() error {
		target := path.Join(h.shimsDir, name, fmt.Sprintf("%s_%s", name, version))
		err := h.defineCommand(ctx, name, version, target, true)
		if err != nil {
			return err
		}

		return h.installCommandBinary(name, version, location, target)
	})
}

func (h *CommandHelper) GetCommand(ctx context.Context, matchers ...q.Matcher) (*model.Command, error) {
	var command model.Command
	err := h.client.Select(matchers...).First(&command)
	if IsQueryNotFound(err) {
		return nil, nil
	}

	return &command, errors.Wrapf(err, "get command failed")
}

func (h *CommandHelper) GetActivatedCommand(ctx context.Context, name string) (*model.Command, error) {
	var command *model.Command

	err := h.client.Atomic(func() error {
		var e error
		command, e = h.GetCommand(ctx, q.Eq("Name", name), q.Eq("Activated", true))
		if e != nil {
			errors.Wrapf(e, "get activated command failed")
		}

		return nil
	})

	return command, err
}

func (h *CommandHelper) GetCommands(ctx context.Context, matchers ...q.Matcher) ([]*model.Command, error) {
	var commands []*model.Command
	err := h.client.Select(matchers...).Find(&commands)
	if IsQueryNotFound(err) {
		return nil, nil
	}

	return commands, err
}

func (h *CommandHelper) activateBinary(ctx context.Context, name, target string) error {
	logger := define.Logger
	fs := define.FS
	binPath := path.Join(h.binDir, name)

	_, err := fs.Stat(binPath)
	if err == nil {
		logger.Debug("remove existed binary", map[string]interface{}{
			"name":   name,
			"target": target,
		})
		fs.Remove(binPath)
	}

	linker := define.GetSymbolLinker()
	err = linker.SymlinkIfPossible(target, binPath)
	if err != nil {
		return errors.Wrapf(err, "create symbol link failed")
	}

	return nil
}

func (h *CommandHelper) Activate(ctx context.Context, name, version string) error {
	return h.client.Atomic(func() error {
		logger := define.Logger

		command, err := h.GetActivatedCommand(ctx, name)
		if err != nil {
			return errors.WithMessagef(err, "deactivate command %s failed", name)
		}

		if command != nil {
			logger.Debug("command found", map[string]interface{}{
				"name":    name,
				"version": command.Version,
			})
			err = h.deactivate(ctx, name)
			if err != nil {
				return err
			}
		}

		logger.Debug("getting command", map[string]interface{}{
			"name":    name,
			"version": version,
		})

		command, err = h.GetCommandByNameAndVersion(ctx, name, version)
		if err != nil {
			return err
		}

		command.Activated = true
		err = h.client.Save(command)
		if err != nil {
			return errors.Wrapf(err, "save command %s failed", name)
		}

		return h.activateBinary(ctx, name, command.Location)
	})
}

func (h *CommandHelper) deactivateBinary(ctx context.Context, name string) error {
	fs := define.FS
	binPath := path.Join(h.binDir, name)

	_, err := fs.Stat(binPath)
	if err != nil {
		return nil
	}

	err = fs.Remove(binPath)
	if err != nil {
		return errors.Wrapf(err, "remove %s failed", name)
	}

	return nil
}

func (h *CommandHelper) deactivate(ctx context.Context, name string) error {
	logger := define.Logger

	command, err := h.GetActivatedCommand(ctx, name)
	if err != nil {
		return err
	}

	if command == nil {
		return errors.Wrapf(ErrCommandNotExists, name)
	}

	logger.Debug("deactivating command", map[string]interface{}{
		"name":    name,
		"version": command.Version,
	})

	command.Activated = false
	err = h.client.Save(command)
	if err != nil {
		return errors.Wrapf(err, "update command %s failed", name)
	}

	return nil
}

func (h *CommandHelper) Deactivate(ctx context.Context, name string) error {
	return h.client.Atomic(func() error {
		err := h.deactivate(ctx, name)
		if err != nil {
			return errors.WithMessagef(err, "deactivate command %s failed", name)
		}

		return h.deactivateBinary(ctx, name)
	})
}

func (h *CommandHelper) Remove(ctx context.Context, name, version string) error {
	return h.client.Atomic(func() error {
		logger := define.Logger
		fs := define.FS

		logger.Debug("remove command", map[string]interface{}{
			"name":    name,
			"version": version,
		})

		command, err := h.GetCommandByNameAndVersion(ctx, name, version)
		if err != nil {
			return err
		}

		err = h.client.DeleteStruct(command)
		if err != nil {
			return errors.Wrapf(err, "delete record %s failed", name)
		}

		if !command.Managed {
			return nil
		}

		err = fs.Remove(command.Location)
		if err != nil {
			return errors.Wrapf(err, "remove command binary %s failed", command.Location)
		}

		return nil
	})
}

func (h *CommandHelper) Upgrade(ctx context.Context, version, path string) (bool, error) {
	name := define.Name
	command, err := h.GetCommandByNameAndVersion(ctx, name, version)
	if err != nil {
		return false, err
	}

	if command != nil {
		return true, nil
	}

	err = h.Install(ctx, name, version, path)
	if err != nil {
		return false, err
	}

	err = h.Activate(ctx, name, version)
	if err != nil {
		return false, err
	}

	return true, nil
}

func NewCommandHelper(client Client) *CommandHelper {
	return &CommandHelper{
		client:   client,
		shimsDir: GetShimsDir(),
		binDir:   GetBinDir(),
	}
}
