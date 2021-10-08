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
	ShimsDir string
	BinDir   string
}

func (h *CommandHelper) GetCommandByNameAndVersion(ctx context.Context, name, version string) (*model.Command, error) {
	return h.GetCommand(ctx, command.Name(name), command.Version(version))
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
	logger := define.Logger

	logger.Debug("checking command", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})

	return utils.WithTx(ctx, h.client, func(client *model.Client) error {
		target := path.Join(h.ShimsDir, name, fmt.Sprintf("%s_%s", name, version))

		_, err := h.client.Command.Create().
			SetName(name).
			SetVersion(version).
			SetLocation(target).
			Save(ctx)

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

func (h *CommandHelper) GetCommands(ctx context.Context, ps ...predicate.Command) ([]*model.Command, error) {
	commands, err := h.client.Command.Query().Where(ps...).All(ctx)
	if model.IsNotFound(err) {
		return nil, nil
	}

	return commands, errors.Wrapf(err, "get commands failed")
}

func (h *CommandHelper) activateBinary(ctx context.Context, name, target string) error {
	fs := define.FS
	binPath := path.Join(h.BinDir, name)

	_, err := fs.Stat(binPath)
	if err == nil {
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

		n, err := h.client.Command.Update().
			Where(command.Name(name), command.Activated(true)).
			SetActivated(false).
			Save(ctx)
		if err != nil {
			return errors.Wrapf(err, "deactivate command failed")
		}

		logger.Debug("activating command", map[string]interface{}{
			"name":        name,
			"version":     version,
			"deactivated": n,
		})
		command, err := h.client.Command.Query().
			Where(command.Name(name), command.Version(version)).
			Only(ctx)

		if err != nil {
			return errors.Wrapf(err, "get command failed")
		}

		command.Activated = true
		err = h.client.Command.UpdateOne(command).Exec(ctx)

		if err != nil {
			return errors.Wrapf(err, "activate command failed")
		}

		return h.activateBinary(ctx, name, command.Location)
	})
}

func NewCommandHelper(client *model.Client) *CommandHelper {
	RootDir := define.Configuration.GetString("cmdr.root")
	return &CommandHelper{
		client:   client,
		ShimsDir: path.Join(RootDir, "shims"),
		BinDir:   path.Join(RootDir, "bin"),
	}
}
