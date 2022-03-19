package manager

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

//go:generate mockgen -destination=mock/storm.go -package=mock github.com/asdine/storm/v3 Query
//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock DBClient

type DBClient interface {
	storm.TypeStore
	Close() error
}

type DatabaseManager struct {
	Client DBClient
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

func (m *DatabaseManager) Define(name string, version string, location string) (core.Command, error) {
	command, _, err := m.getOrNew(name, version)
	if err != nil {
		return nil, errors.Wrapf(err, "define command failed")
	}

	command.Location = location
	core.GetLogger().Debug("defining command", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})

	err = m.Client.Save(command)
	if err != nil {
		return nil, errors.Wrapf(err, "save command failed")
	}

	return command, nil
}

func (m *DatabaseManager) Undefine(name string, version string) error {
	command, found, err := m.getOrNew(name, version)
	if err != nil {
		return errors.Wrapf(err, "undefine command failed")
	}

	if !found {
		return nil
	}

	core.GetLogger().Debug("undefining command", map[string]interface{}{
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

	core.GetLogger().Debug("activating command", map[string]interface{}{
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

	core.GetLogger().Debug("deactivating commands", map[string]interface{}{
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

type DatabaseMigrator struct {
	dbFactory func() (DBClient, error)
}

func (m *DatabaseMigrator) Init() error {
	logger := core.GetLogger()

	db, err := m.dbFactory()
	if err != nil {
		return errors.Wrapf(err, "open database failed")
	}
	defer utils.CallClose(db)

	for name, model := range map[string]interface{}{
		"command": &Command{},
	} {
		logger.Debug("initializing database model", map[string]interface{}{
			"model": name,
		})
		err := db.Init(model)
		if err != nil {
			return errors.Wrapf(err, "init database failed")
		}

		logger.Debug("indexing database model", map[string]interface{}{
			"model": name,
		})
		err = db.ReIndex(model)
		if err != nil {
			return errors.Wrapf(err, "reindex database failed")
		}
	}

	return nil
}

func NewDatabaseMigrator(dbFactory func() (DBClient, error)) *DatabaseMigrator {
	return &DatabaseMigrator{
		dbFactory: dbFactory,
	}
}

func init() {
	var (
		_ core.Command        = (*Command)(nil)
		_ core.CommandQuery   = (*CommandQuery)(nil)
		_ core.CommandManager = (*DatabaseManager)(nil)
	)

	core.RegisterCommandManagerFactory(core.CommandProviderDatabase, func(cfg core.Configuration) (core.CommandManager, error) {
		dbPath := cfg.GetString(core.CfgKeyCmdrDatabasePath)

		db, err := storm.Open(dbPath)
		if err != nil {
			return nil, errors.Wrapf(err, "open database failed")
		}

		return NewDatabaseManager(db), nil
	})

	core.RegisterInitializerFactory("database-migrator", func(cfg core.Configuration) (core.Initializer, error) {
		return NewDatabaseMigrator(func() (DBClient, error) {
			return storm.Open(cfg.GetString(core.CfgKeyCmdrDatabasePath))
		}), nil
	})
}
