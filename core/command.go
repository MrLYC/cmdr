package core

import (
	"context"
	"fmt"
	"path"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/model/command"
	"github.com/mrlyc/cmdr/model/predicate"
	"github.com/mrlyc/cmdr/utils"
)

type CommandHelper struct {
	client   *model.Client
	shimsDir string
	binDir   string
}

func (h *CommandHelper) GetCommandByNameAndVersion(ctx context.Context, name, version string) (*model.Command, error) {
	return h.GetCommand(ctx, command.Name(name), command.Version(version))
}

func (h *CommandHelper) defineCommand(ctx context.Context, name, version, location string, managed bool) error {
	logger := define.Logger

	logger.Debug("checking command", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})

	h.client.Command.Create().
		SetName(name).
		SetVersion(version).
		SetLocation(location).
		SetManaged(managed).
		SaveX(ctx)

	return nil
}

func (h *CommandHelper) Define(ctx context.Context, name, version, location string) error {
	return utils.WithTx(ctx, h.client, func(client *model.Client) error {
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
	return utils.WithTx(ctx, h.client, func(client *model.Client) error {
		target := path.Join(h.shimsDir, name, fmt.Sprintf("%s_%s", name, version))
		err := h.defineCommand(ctx, name, version, target, true)
		if err != nil {
			return err
		}

		return h.installCommandBinary(name, version, location, target)
	})
}

func (h *CommandHelper) GetCommand(ctx context.Context, ps ...predicate.Command) (*model.Command, error) {
	command, err := h.client.Command.Query().Where(ps...).Only(ctx)
	if model.IsNotFound(err) {
		return nil, nil
	}

	return command, errors.Wrapf(err, "get command failed")
}

func (h *CommandHelper) GetActivatedCommand(ctx context.Context, name string) (*model.Command, error) {
	command, err := h.client.Command.Query().Where(command.Name(name), command.Activated(true)).
		Only(ctx)
	if model.IsNotFound(err) {
		return nil, nil
	}

	return command, errors.Wrapf(err, "get activated command failed")
}

func (h *CommandHelper) GetCommands(ctx context.Context, ps ...predicate.Command) ([]*model.Command, error) {
	commands, err := h.client.Command.Query().Where(ps...).All(ctx)
	if model.IsNotFound(err) {
		return nil, nil
	}

	return commands, errors.Wrapf(err, "get commands failed")
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
	return utils.WithTx(ctx, h.client, func(client *model.Client) error {
		logger := define.Logger

		n := h.client.Command.Update().
			Where(command.Name(name), command.Activated(true)).
			SetActivated(false).
			SaveX(ctx)

		logger.Debug("activating command", map[string]interface{}{
			"name":        name,
			"version":     version,
			"deactivated": n,
		})

		n = h.client.Command.Update().
			Where(command.Name(name), command.Version(version)).
			SetActivated(true).
			SaveX(ctx)

		cmd, err := h.GetCommandByNameAndVersion(ctx, name, version)
		if err != nil {
			return err
		}
		return h.activateBinary(ctx, name, cmd.Location)
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

func (h *CommandHelper) Deactivate(ctx context.Context, name string) error {
	return utils.WithTx(ctx, h.client, func(client *model.Client) error {
		logger := define.Logger

		n := h.client.Command.Update().
			Where(command.Name(name)).
			SetActivated(false).
			SaveX(ctx)

		logger.Debug("deactivating command", map[string]interface{}{
			"name":        name,
			"deactivated": n,
		})

		return h.deactivateBinary(ctx, name)
	})
}

func (h *CommandHelper) Remove(ctx context.Context, name, version string) error {
	return utils.WithTx(ctx, h.client, func(client *model.Client) error {
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

		h.client.Command.DeleteOneID(command.ID).ExecX(ctx)
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

func NewCommandHelper(client *model.Client) *CommandHelper {
	return &CommandHelper{
		client:   client,
		shimsDir: GetShimsDir(),
		binDir:   GetBinDir(),
	}
}
