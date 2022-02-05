package operator

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
)

type CommandsQuerier struct {
	BaseOperator
	matchers []q.Matcher
}

func (c *CommandsQuerier) String() string {
	return "commands-querier"
}

func (c *CommandsQuerier) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	var commands []*model.Command
	client := GetDBClientFromContext(ctx)

	err := client.Select(c.matchers...).Find(&commands)
	switch errors.Cause(err) {
	case nil, storm.ErrNotFound:
	default:
		return ctx, errors.Wrapf(err, "query commands failed")
	}

	logger.Debug("commands queried", map[string]interface{}{
		"count": len(commands),
	})

	return context.WithValue(ctx, define.ContextKeyCommands, commands), nil
}

func NewCommandsQuerier(matchers []q.Matcher) *CommandsQuerier {
	return &CommandsQuerier{
		matchers: matchers,
	}
}

func NewFullCommandsQuerier(name, version, location string, activated bool) *CommandsQuerier {
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

func NewSimpleCommandsQuerier(name, version string) *CommandsQuerier {
	return NewCommandsQuerier([]q.Matcher{q.Eq("Name", name), q.Eq("Version", version)})
}

func NewNamedCommandsQuerier(name string) *CommandsQuerier {
	return NewCommandsQuerier([]q.Matcher{q.Eq("Name", name)})
}
