package core

import (
	"context"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
)

type CommandHelper struct {
	client *model.Client
}

func (h *CommandHelper) Install(ctx context.Context, name, version, location string) error {
	logger := define.Logger

	logger.Debug("installing command", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})
	_, err := h.client.Command.Create().
		SetName(name).
		SetVersion(version).
		SetLocation(location).
		Save(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (h *CommandHelper) Query() *model.CommandQuery {
	return h.client.Command.Query()
}

func NewCommandHelper(client *model.Client) *CommandHelper {
	return &CommandHelper{
		client: client,
	}
}
