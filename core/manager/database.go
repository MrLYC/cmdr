package manager

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type DatabaseManager struct {
	Client  core.Database
	manager core.CommandManager
}

func (m *DatabaseManager) Close() error {
	err := m.manager.Close()
	if err != nil {
		return errors.Wrapf(err, "close %v manager failed", m.manager.Provider())
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
	defined, err := m.manager.Define(name, version, location)
	if err != nil {
		return nil, err
	}

	location = defined.GetLocation()

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

	return m.manager.Undefine(name, version)
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

	return m.manager.Activate(name, version)
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

	return m.manager.Deactivate(name)
}

func NewDatabaseManager(db core.Database, manager core.CommandManager) *DatabaseManager {
	return &DatabaseManager{
		Client:  db,
		manager: manager,
	}
}

func init() {
	core.RegisterCommandManagerFactory(core.CommandProviderDatabase, func(cfg core.Configuration) (core.CommandManager, error) {
		mgr, err := core.NewCommandManager(core.CommandProviderBinary, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "new manager binary failed")
		}

		db, err := core.GetDatabase()
		if err != nil {
			return nil, errors.Wrapf(err, "open database failed")
		}

		return NewDatabaseManager(db, mgr), nil
	})
}
