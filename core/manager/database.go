package manager

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

//go:generate mockgen -destination=mock/storm.go -package=mock github.com/asdine/storm/v3 Query
//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock DBClient

type DBClient interface {
	storm.TypeStore
	Close() error
}

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

type DatabaseManager struct {
	Client DBClient
}

func (m *DatabaseManager) Init() error {
	var command Command
	logger := core.Logger

	logger.Debug("initializing database model", map[string]interface{}{
		"model": "command",
	})
	err := m.Client.Init(&command)
	if err != nil {
		return errors.Wrapf(err, "init database failed")
	}

	logger.Debug("indexing database model", map[string]interface{}{
		"model": "command",
	})
	err = m.Client.ReIndex(&command)
	if err != nil {
		return errors.Wrapf(err, "reindex database failed")
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

func (m *DatabaseManager) Provider() core.CommandProvider {
	return core.CommandProviderDatabase
}

func (m *DatabaseManager) Query() (core.CommandQuery, error) {
	return NewCommandQuery(m.Client), nil
}

func (m *DatabaseManager) getOrNew(name string, version string) (*Command, bool, error) {
	var found bool
	var command Command
	err := m.Client.Select(q.Eq("Name", name), q.Eq("Version", version)).First(&command)
	switch errors.Cause(err) {
	case nil:
		found = true
	case storm.ErrNotFound:
		found = false
	default:
		return nil, false, errors.Wrapf(err, "get command failed")
	}

	command.Name = name
	command.Version = version

	return &command, found, nil
}

func (m *DatabaseManager) Define(name string, version string, location string) error {
	command, _, err := m.getOrNew(name, version)
	if err != nil {
		return errors.Wrapf(err, "define command failed")
	}

	command.Location = location
	core.Logger.Debug("defining command", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})

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

	core.Logger.Debug("undefining command", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	if command.Activated {
		return errors.Wrapf(core.ErrCommandAlreadyActivated, "command %s:%s is activated", name, version)
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

	core.Logger.Debug("activating command", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	err = m.Deactivate(name)
	if err != nil {
		return errors.Wrapf(err, "deactivate commands failed")
	}

	command.Activated = true

	err = m.Client.Save(command)
	if err != nil {
		return errors.Wrapf(err, "save command failed")
	}

	return nil
}

func (m *DatabaseManager) Deactivate(name string) error {
	var commands []*Command
	err := m.Client.Select(
		q.Eq("Name", name),
		q.Eq("Activated", true),
	).Find(&commands)
	switch errors.Cause(err) {
	case nil:
	case storm.ErrNotFound:
		return nil
	default:
		return errors.Wrapf(err, "select commands failed")
	}

	core.Logger.Debug("deactivating commands", map[string]interface{}{
		"name": name,
	})

	for _, cmd := range commands {
		cmd.Activated = false
		err := m.Client.Save(cmd)
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

func newDatabaseManagerByConfiguration(cfg core.Configuration) (*DatabaseManager, error) {
	dbPath := cfg.GetString(core.CfgKeyCmdrDatabasePath)

	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, errors.Wrapf(err, "open database failed")
	}

	return NewDatabaseManager(db), nil
}

func init() {
	var (
		_ core.Command        = (*Command)(nil)
		_ core.CommandQuery   = (*CommandQuery)(nil)
		_ core.CommandManager = (*DatabaseManager)(nil)
		_ core.Initializer    = (*DatabaseManager)(nil)
	)

	core.RegisterCommandManagerFactory(core.CommandProviderDatabase, func(cfg core.Configuration) (core.CommandManager, error) {
		return newDatabaseManagerByConfiguration(cfg)
	})

	core.RegisterInitializerFactory("database", func(cfg core.Configuration) (core.Initializer, error) {
		return newDatabaseManagerByConfiguration(cfg)
	})
}
