package manager

import (
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/cmdr"
)

//go:generate mockgen -destination=mock/storm.go -package=mock github.com/asdine/storm/v3 Query
//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock DBClient

type DBClient interface {
	storm.TypeStore
	Close() error
}

var ErrCommandAlreadyActivated = fmt.Errorf("command already activated")

type Command struct {
	IDField        int    `storm:"id,increment"`
	NameField      string `storm:"name,index" json:"name"`
	VersionField   string `storm:"version,index" json:"version"`
	ActivatedField bool   `storm:"activated,index" json:"activated"`
	LocationField  string `storm:"location" json:"location"`
}

func (c *Command) Name() string {
	return c.NameField
}

func (c *Command) Version() string {
	return c.VersionField
}

func (c *Command) Activated() bool {
	return c.ActivatedField
}

func (c *Command) Location() string {
	return c.LocationField
}

type CommandQuery struct {
	Client   storm.TypeStore
	matchers []q.Matcher
	query    storm.Query
}

func (c *CommandQuery) WithName(name string) cmdr.CommandQuery {
	c.matchers = append(c.matchers, q.Eq("NameField", name))
	return c
}

func (c *CommandQuery) WithVersion(version string) cmdr.CommandQuery {
	c.matchers = append(c.matchers, q.Eq("VersionField", version))
	return c
}

func (c *CommandQuery) WithActivated(activated bool) cmdr.CommandQuery {
	c.matchers = append(c.matchers, q.Eq("ActivatedField", activated))
	return c
}

func (c *CommandQuery) WithLocation(location string) cmdr.CommandQuery {
	c.matchers = append(c.matchers, q.Eq("LocationField", location))
	return c
}

func (c *CommandQuery) Done() storm.Query {
	if c.query == nil {
		c.query = c.Client.Select(c.matchers...)
	}
	return c.query
}

func (c *CommandQuery) All() ([]cmdr.Command, error) {
	var commands []*Command
	err := c.Done().Find(&commands)
	if err != nil {
		return nil, err
	}

	result := make([]cmdr.Command, 0, len(commands))
	for _, cmd := range commands {
		result = append(result, cmd)
	}
	return result, nil
}

func (c *CommandQuery) One() (cmdr.Command, error) {
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

type DatabaseManager struct {
	Client DBClient
}

func (m *DatabaseManager) Init() error {
	var command Command

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

func (m *DatabaseManager) Close() error {
	err := m.Client.Close()
	if err != nil {
		return errors.Wrapf(err, "close database failed")
	}

	return nil
}

func (m *DatabaseManager) Provider() cmdr.CommandProvider {
	return cmdr.CommandProviderDatabase
}

func (m *DatabaseManager) Query() (cmdr.CommandQuery, error) {
	return NewCommandQuery(m.Client), nil
}

func (m *DatabaseManager) getOrNew(name string, version string) (*Command, bool, error) {
	var found bool
	var command Command
	err := m.Client.Select(q.Eq("NameField", name), q.Eq("VersionField", version)).First(&command)
	switch errors.Cause(err) {
	case nil:
		found = true
	case storm.ErrNotFound:
		found = false
	default:
		return nil, false, errors.Wrapf(err, "get command failed")
	}

	command.NameField = name
	command.VersionField = version

	return &command, found, nil
}

func (m *DatabaseManager) Define(name string, version string, location string) error {
	command, _, err := m.getOrNew(name, version)
	if err != nil {
		return errors.Wrapf(err, "define command failed")
	}

	command.LocationField = location

	err = m.Client.Save(command)
	if err != nil {
		return errors.Wrapf(err, "save command failed")
	}

	return nil
}

func (m *DatabaseManager) Undefine(name string, version string) error {
	command, found, err := m.getOrNew(name, version)
	if err != nil {
		return errors.Wrapf(err, "undefine command failed")
	}

	if !found {
		return nil
	}

	if command.ActivatedField {
		return errors.Wrapf(ErrCommandAlreadyActivated, "command %s:%s is activated", name, version)
	}

	err = m.Client.DeleteStruct(command)
	if err != nil {
		return errors.Wrapf(err, "delete command failed")
	}

	return nil
}

func (m *DatabaseManager) Activate(name string, version string) error {
	command, found, err := m.getOrNew(name, version)
	if err != nil {
		return errors.Wrapf(err, "activate command failed")
	}

	if !found {
		return errors.Errorf("command %s(%s) not found", name, version)
	}

	err = m.Deactivate(name)
	if err != nil {
		return errors.Wrapf(err, "deactivate commands failed")
	}

	command.ActivatedField = true

	err = m.Client.Save(command)
	if err != nil {
		return errors.Wrapf(err, "save command failed")
	}

	return nil
}

func (m *DatabaseManager) Deactivate(name string) error {
	var commands []*Command
	err := m.Client.Select(
		q.Eq("NameField", name),
		q.Eq("ActivatedField", true),
	).Find(&commands)
	switch errors.Cause(err) {
	case nil:
	case storm.ErrNotFound:
		return nil
	default:
		return errors.Wrapf(err, "select commands failed")
	}

	for _, cmd := range commands {
		cmd.ActivatedField = false
		err := m.Client.Update(cmd)
		if err != nil {
			return errors.Wrapf(err, "deactivate command failed")
		}
	}

	return nil
}

func NewDatabaseManager(db DBClient) *DatabaseManager {
	return &DatabaseManager{
		Client: db,
	}
}

func init() {
	var (
		_ cmdr.Command        = (*Command)(nil)
		_ cmdr.CommandQuery   = (*CommandQuery)(nil)
		_ cmdr.CommandManager = (*DatabaseManager)(nil)
	)

	cmdr.RegisterCommandManagerFactory(cmdr.CommandProviderDatabase, func(cfg cmdr.Configuration, opts ...cmdr.Option) (cmdr.CommandManager, error) {
		dbPath := cfg.GetString(cmdr.CfgKeyCmdrDatabasePath)

		db, err := storm.Open(dbPath)
		if err != nil {
			return nil, errors.Wrapf(err, "open database failed")
		}

		return NewDatabaseManager(db), nil
	})
}
