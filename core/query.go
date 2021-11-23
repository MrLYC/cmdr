package core

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
)

type CommandListQuerier struct {
	BaseStep
	matchers []q.Matcher
}

func (c *CommandListQuerier) String() string {
	return "commands-querier"
}

func (c *CommandListQuerier) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	var commands []*model.Command
	client := GetDBClientFromContext(ctx)
	err := client.Select(c.matchers...).Find(&commands)
	if errors.Cause(err) == storm.ErrNotFound {
		return ctx, nil
	} else if err != nil {
		return ctx, errors.Wrap(err, "query command failed")
	}

	logger.Debug("commands queried", map[string]interface{}{
		"count": len(commands),
	})

	return context.WithValue(ctx, define.ContextKeyCommands, commands), nil
}

func NewCommandsQuerier(matchers []q.Matcher) *CommandListQuerier {
	return &CommandListQuerier{
		matchers: matchers,
	}
}

func NewSimpleCommandsQuerier(name, version, location string, activated bool) *CommandListQuerier {
	logger := define.Logger
	filters := make([]q.Matcher, 0)

	if name != "" {
		logger.Debug("filter by name", map[string]interface{}{
			"name": name,
		})
		filters = append(filters, q.Eq("Name", name))
	}

	if version != "" {
		logger.Debug("filter by version", map[string]interface{}{
			"version": version,
		})
		filters = append(filters, q.Eq("Version", version))
	}

	if location != "" {
		logger.Debug("filter by location", map[string]interface{}{
			"location": location,
		})
		filters = append(filters, q.Eq("Location", location))
	}

	if activated {
		logger.Debug("filter by activated", map[string]interface{}{
			"activated": activated,
		})
		filters = append(filters, q.Eq("Activated", activated))
	}

	return NewCommandsQuerier(filters)
}

type CommandQuerier struct {
	BaseStep
	matchers []q.Matcher
}

func (c *CommandQuerier) String() string {
	return "command-querier"
}

func (c *CommandQuerier) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	var command model.Command
	client := GetDBClientFromContext(ctx)
	err := client.Select(c.matchers...).First(&command)
	if errors.Cause(err) == storm.ErrNotFound {
		return ctx, nil
	} else if err != nil {
		return ctx, errors.Wrap(err, "query command failed")
	}

	logger.Info("command queried", map[string]interface{}{
		"name":    command.Name,
		"version": command.Version,
	})
	return context.WithValue(ctx, define.ContextKeyCommand, &command), nil
}

func NewCommandQuerier(matchers []q.Matcher) *CommandQuerier {
	return &CommandQuerier{
		matchers: matchers,
	}
}

func NewCommandListQuerierByNameAndVersion(name, version string) *CommandListQuerier {
	return NewCommandsQuerier(
		[]q.Matcher{q.Eq("Name", name), q.Eq("Version", version)},
	)
}
