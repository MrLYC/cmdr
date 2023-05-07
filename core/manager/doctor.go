package manager

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type DoctorManager struct {
	binaryMgr   core.CommandManager
	databaseMgr core.CommandManager
}

func (m *DoctorManager) each(fn func(mgr core.CommandManager) error) error {
	for _, mgr := range []core.CommandManager{m.binaryMgr, m.databaseMgr} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *DoctorManager) reverseEach(fn func(mgr core.CommandManager) error) error {
	for _, mgr := range []core.CommandManager{m.databaseMgr, m.binaryMgr} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *DoctorManager) all(fn func(mgr core.CommandManager) error) error {
	var errs error
	for _, mgr := range []core.CommandManager{m.binaryMgr, m.databaseMgr} {
		err := fn(mgr)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (d *DoctorManager) Define(name, version, location string) (core.Command, error) {
	var command core.Command
	err := d.all(func(mgr core.CommandManager) error {
		cmd, err := mgr.Define(name, version, location)
		if err != nil {
			return err
		}
		command = cmd
		return nil
	})
	return command, err
}

func (d *DoctorManager) Undefine(name, version string) error {
	return d.all(func(mgr core.CommandManager) error {
		return mgr.Undefine(name, version)
	})
}

func (d *DoctorManager) Activate(name, version string) error {
	return d.all(func(mgr core.CommandManager) error {
		return mgr.Activate(name, version)
	})
}

func (d *DoctorManager) Deactivate(name string) error {
	return d.all(func(mgr core.CommandManager) error {
		return mgr.Deactivate(name)
	})
}

func (d *DoctorManager) Close() error {
	return d.all(func(mgr core.CommandManager) error {
		return mgr.Close()
	})
}

func (d *DoctorManager) Provider() core.CommandProvider {
	return core.CommandProviderDoctor
}

func (d *DoctorManager) Query() (core.CommandQuery, error) {
	mainQuery, mainErr := d.binaryMgr.Query()
	if mainErr != nil {
		return d.databaseMgr.Query()
	}

	recorderQuery, recorderErr := d.databaseMgr.Query()
	if recorderErr != nil {
		return mainQuery, mainErr
	}

	var queriedCommands []core.Command

	// merge two queries
	commands, mainErr := mainQuery.All()
	if mainErr == nil {
		queriedCommands = append(queriedCommands, commands...)
	}

	recorderCommands, recorderErr := recorderQuery.All()
	if recorderErr == nil {
		queriedCommands = append(queriedCommands, recorderCommands...)
	}

	indexes := make(map[string]int, len(queriedCommands))
	merged := make([]*Command, 0, len(queriedCommands))
	for i, cmd := range queriedCommands {
		name := cmd.GetName()
		version := cmd.GetVersion()

		key := fmt.Sprintf("%s-%s", name, version)
		index, ok := indexes[key]
		if ok {
			// update by recorder
			merged[index].Activated = cmd.GetActivated()
			continue
		}

		indexes[key] = i
		merged = append(merged, &Command{
			Name:      name,
			Version:   version,
			Activated: cmd.GetActivated(),
			Location:  cmd.GetLocation(),
		})
	}

	return NewCommandFilter(merged), nil
}

func NewDoctorManager(binaryMgr core.CommandManager, databaseMgr core.CommandManager) *DoctorManager {
	return &DoctorManager{
		binaryMgr:   binaryMgr,
		databaseMgr: databaseMgr,
	}
}

type CommandDoctor struct {
	core.CommandManager
}

func (d *CommandDoctor) Fix() error {
	logger := core.GetLogger()

	query, err := d.Query()
	if err != nil {
		return errors.Wrapf(err, "make query failed")
	}

	commands, err := query.All()
	if err != nil {
		return errors.Wrapf(err, "query commands failed")
	}

	var availableCommands []core.Command
	for _, cmd := range commands {
		name := cmd.GetName()
		version := cmd.GetVersion()
		location := cmd.GetLocation()
		activated := cmd.GetActivated()

		logger.Debug("checking command", map[string]interface{}{
			"name":     name,
			"version":  version,
			"location": location,
		})

		_, err := os.Stat(location)
		if err == nil {
			logger.Debug("command is available", map[string]interface{}{
				"name":    name,
				"version": version,
			})
			availableCommands = append(availableCommands, cmd)
			continue
		}

		logger.Warn("command is not available", map[string]interface{}{
			"name":    name,
			"version": version,
		})

		if activated {
			logger.Info("deactivating command", map[string]interface{}{
				"name": name,
			})
			err = d.Deactivate(name)
			if err != nil {
				logger.Warn("deactivate command failed, try to remove it", map[string]interface{}{
					"name": name,
				})
			}
		}

		logger.Info("removing command", map[string]interface{}{
			"name":    name,
			"version": version,
		})
		err = d.Undefine(name, version)
		if err != nil {
			logger.Error("remove command failed, aborted", map[string]interface{}{
				"name":    name,
				"version": version,
			})
		}
	}

	for _, cmd := range availableCommands {
		name := cmd.GetName()
		version := cmd.GetVersion()
		location := cmd.GetLocation()
		activated := cmd.GetActivated()

		_, err := d.Define(name, version, location)
		if err != nil {
			logger.Warn("re-define command failed, continue", map[string]interface{}{
				"name":     name,
				"version":  version,
				"location": location,
			})
		}

		if activated {
			err = d.Activate(name, version)
			if err != nil {
				logger.Warn("re-activate command failed, continue", map[string]interface{}{
					"name":    name,
					"version": version,
				})
			}
		}
	}

	return nil
}

func NewCommandDoctor(manager core.CommandManager) *CommandDoctor {
	return &CommandDoctor{
		CommandManager: manager,
	}
}

func init() {
	var _ core.CommandManager = (*DoctorManager)(nil)

	core.RegisterCommandManagerFactory(core.CommandProviderDoctor, func(cfg core.Configuration) (core.CommandManager, error) {
		mainMgr, err := core.NewCommandManager(core.CommandProviderBinary, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "new main command manager failed")
		}

		recorderMgr, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "new recorder command manager failed")
		}

		return NewDoctorManager(mainMgr, recorderMgr), nil
	})
}
