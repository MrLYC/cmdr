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

func (h *CommandHelper) installCommandBinary(name, version, location string) (string, error) {
	fs := define.FS
	logger := define.Logger

	shimsDir := path.Join(h.ShimsDir, name)
	logger.Debug("creating command shims dir", map[string]interface{}{
		"name": name,
		"dir":  shimsDir,
	})
	utils.CheckError(fs.MkdirAll(shimsDir, 0755))

	target := path.Join(shimsDir, fmt.Sprintf("%s_%s", name, version))
	logger.Debug("coping command", map[string]interface{}{
		"name":     name,
		"location": location,
		"target":   target,
	})
	err := utils.CopyFile(location, target)
	if err != nil {
		return "", errors.WithMessagef(err, "install command %s failed", target)
	}

	err = fs.Chmod(target, 0755)
	if err != nil {
		return "", errors.WithMessagef(err, "change command mode %s failed", target)
	}

	return target, nil
}

func (h *CommandHelper) Install(ctx context.Context, name, version, location string) error {
	logger := define.Logger

	logger.Debug("checking command", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})
	command, err := h.GetCommandByNameAndVersion(ctx, name, version)
	if err != nil {
		return err
	}

	if command != nil {
		return errors.Wrapf(ErrCommandAlreadyExists, "name %s, version %s", name, version)
	}

	target, err := h.installCommandBinary(name, version, location)
	if err != nil {
		return err
	}

	_, err = h.client.Command.Create().
		SetName(name).
		SetVersion(version).
		SetLocation(target).
		Save(ctx)

	if err != nil {
		return err
	}

	return nil
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

func NewCommandHelper(client *model.Client) *CommandHelper {
	RootDir := define.Configuration.GetString("cmdr.root")
	return &CommandHelper{
		client:   client,
		ShimsDir: path.Join(RootDir, "shims"),
		BinDir:   path.Join(RootDir, "bin"),
	}
}
