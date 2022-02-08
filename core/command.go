package core

import (
	"path/filepath"

	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core/model"
	"github.com/mrlyc/cmdr/define"
)

type CommandQuery struct {
	Client   define.DBClient
	matchers []q.Matcher
}

func (c *CommandQuery) WithName(name string) *CommandQuery {
	c.matchers = append(c.matchers, q.Eq("Name", name))
	return c
}

func (c *CommandQuery) WithVersion(version string) *CommandQuery {
	c.matchers = append(c.matchers, q.Eq("Version", version))
	return c
}

func (c *CommandQuery) WithActivated(activated bool) *CommandQuery {
	c.matchers = append(c.matchers, q.Eq("Activated", activated))
	return c
}

func (c *CommandQuery) WithLocation(location string) *CommandQuery {
	c.matchers = append(c.matchers, q.Eq("Location", location))
	return c
}

func (c *CommandQuery) Done() define.DBQuery {
	return c.Client.Select(c.matchers...)
}

func NewCommandQuery(db define.DBClient) *CommandQuery {
	return &CommandQuery{
		Client: db,
	}
}

type CommandManager struct {
	Client define.DBClient
}

func (m *CommandManager) Delete(cmd *model.Command) error {
	cmd.Activated = true
	err := m.Client.DeleteStruct(cmd)
	if err != nil {
		return errors.Wrapf(err, "delete command failed")
	}
	return nil
}

func (m *CommandManager) Activate(cmd *model.Command) error {
	cmd.Activated = true
	err := m.Client.Update(cmd)
	if err != nil {
		return errors.Wrapf(err, "activate command failed")
	}
	return nil
}

func (m *CommandManager) Deactivate(cmd *model.Command) error {
	cmd.Activated = false
	err := m.Client.Update(cmd)
	if err != nil {
		return errors.Wrapf(err, "deactivate command failed")
	}
	return nil
}

func (m *CommandManager) DeactivateAll(name string) error {
	var commands []model.Command
	err := m.Client.Select(q.Eq("Name", name)).Find(&commands)
	if err != nil {
		return errors.Wrapf(err, "find commands failed")
	}

	for _, cmd := range commands {
		err = m.Deactivate(&cmd)
		if err != nil {
			return errors.Wrapf(err, "deactivate command failed")
		}
	}

	return nil
}

func (m *CommandManager) Query() *CommandQuery {
	return NewCommandQuery(m.Client)
}

func (m *CommandManager) Init() error {
	var command model.Command

	err := m.Client.Init(&command)
	if err != nil {
		return errors.Wrapf(err, "init command failed")
	}

	err = m.Client.ReIndex(&command)
	if err != nil {
		return errors.Wrapf(err, "reindex command failed")
	}

	return nil
}

func NewCommandManagerWith(db define.DBClient) *CommandManager {
	return &CommandManager{
		Client: db,
	}
}

func NewCommandManager(root string) (*CommandManager, error) {
	Client, err := NewDBClient(filepath.Join(root, "cmdr.db"))
	if err != nil {
		return nil, errors.Wrapf(err, "create db client failed")
	}

	return NewCommandManagerWith(Client), nil
}
