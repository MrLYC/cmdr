package manager

import (
	. "github.com/ahmetb/go-linq/v3"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type Command struct {
	ID        int    `storm:"increment"`
	Name      string `storm:"index" json:"name"`
	Version   string `storm:"index" json:"version"`
	Activated bool   `storm:"index" json:"activated"`
	Location  string `storm:"" json:"location"`
}

func (c *Command) GetName() string {
	return c.Name
}

func (c *Command) GetVersion() string {
	return c.Version
}

func (c *Command) GetActivated() bool {
	return c.Activated
}

func (c *Command) GetLocation() string {
	return c.Location
}

type CommandFilter struct {
	commands []*Command
}

func (f *CommandFilter) Filter(fn func(b interface{}) bool) *CommandFilter {
	From(f.commands).Where(fn).ToSlice(&f.commands)
	return f
}

func (f *CommandFilter) WithName(name string) core.CommandQuery {
	return f.Filter(func(b interface{}) bool {
		return b.(*Command).Name == name
	})
}

func (f *CommandFilter) WithVersion(version string) core.CommandQuery {
	return f.Filter(func(b interface{}) bool {
		return b.(*Command).Version == version
	})
}

func (f *CommandFilter) WithActivated(activated bool) core.CommandQuery {
	return f.Filter(func(b interface{}) bool {
		return b.(*Command).Activated == activated
	})
}

func (f *CommandFilter) WithLocation(location string) core.CommandQuery {
	return f.Filter(func(b interface{}) bool {
		return b.(*Command).Location == location
	})
}
func (f *CommandFilter) All() ([]core.Command, error) {
	commands := make([]core.Command, 0, len(f.commands))
	for _, b := range f.commands {
		commands = append(commands, b)
	}

	return commands, nil
}

func (f *CommandFilter) One() (core.Command, error) {
	if len(f.commands) == 0 {
		return nil, errors.Wrapf(core.ErrBinaryNotFound, "commands not found")
	}

	return f.commands[0], nil
}

func (f *CommandFilter) Count() (int, error) {
	return len(f.commands), nil
}

func (f *CommandFilter) AddCommand(commands ...core.Command) {
	for _, command := range commands {
		f.commands = append(f.commands, &Command{
			Name:      command.GetName(),
			Version:   command.GetVersion(),
			Activated: command.GetActivated(),
			Location:  command.GetLocation(),
		})
	}
}

func NewCommandFilter(commands []*Command) *CommandFilter {
	return &CommandFilter{commands}
}

type CommandQuery struct {
	Client   storm.TypeStore
	matchers []q.Matcher
	query    storm.Query
}

func (c *CommandQuery) WithName(name string) core.CommandQuery {
	c.matchers = append(c.matchers, q.Eq("Name", name))
	return c
}

func (c *CommandQuery) WithVersion(version string) core.CommandQuery {
	c.matchers = append(c.matchers, q.Eq("Version", version))
	return c
}

func (c *CommandQuery) WithActivated(activated bool) core.CommandQuery {
	c.matchers = append(c.matchers, q.Eq("Activated", activated))
	return c
}

func (c *CommandQuery) WithLocation(location string) core.CommandQuery {
	c.matchers = append(c.matchers, q.Eq("Location", location))
	return c
}

func (c *CommandQuery) Done() storm.Query {
	if c.query == nil {
		c.query = c.Client.Select(c.matchers...)
	}
	return c.query
}

func (c *CommandQuery) All() ([]core.Command, error) {
	var commands []*Command
	err := c.Done().Find(&commands)
	if err != nil {
		return nil, err
	}

	result := make([]core.Command, 0, len(commands))
	for _, cmd := range commands {
		result = append(result, cmd)
	}
	return result, nil
}

func (c *CommandQuery) One() (core.Command, error) {
	var cmd Command
	err := c.Done().First(&cmd)
	if err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (c *CommandQuery) Count() (int, error) {
	var cmd Command
	return c.Done().Count(&cmd)
}

func NewCommandQuery(db storm.TypeStore) *CommandQuery {
	return &CommandQuery{
		Client: db,
	}
}

func init() {
	core.RegisterDatabaseModel(core.ModelTypeCommand, &Command{})
}
