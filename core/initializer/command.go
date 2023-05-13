package initializer

import (
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type CmdrUpdater struct {
	name     string
	version  string
	location string
	manager  core.CommandManager
}

func (c *CmdrUpdater) getActivatedCmdrVersion() string {
	query, err := c.manager.Query()
	if err != nil {
		return ""
	}

	command, err := query.
		WithName(c.name).
		WithActivated(true).
		One()

	if err != nil {
		return ""
	}

	return command.GetVersion()
}

func (c *CmdrUpdater) removeLegacies(safeVersions []string) error {
	logger := core.GetLogger()

	query, err := c.manager.Query()
	if err != nil {
		return errors.Wrapf(err, "failed to create command query")
	}

	commands, err := query.
		WithName(c.name).
		All()

	if err != nil {
		return errors.Wrapf(err, "failed to get commands")
	}

	var errs error
	for _, command := range commands {
		logger.Debug("checking legacy command", map[string]interface{}{
			"command": command,
		})

		if command.GetActivated() {
			continue
		}

		version := command.GetVersion()
		isSafe := false
		for _, safeVersion := range safeVersions {
			if safeVersion == version {
				isSafe = true
				break
			}
		}

		if isSafe {
			continue
		}

		logger.Info("removing legacy cmdr", map[string]interface{}{
			"command": command,
		})
		err = c.manager.Undefine(c.name, version)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (c *CmdrUpdater) Init(isUpgrade bool) error {
	logger := core.GetLogger()
	safeVersion := []string{c.version}
	version := c.getActivatedCmdrVersion()
	if version != "" {
		safeVersion = append(safeVersion, version)
	}

	logger.Debug("update command", map[string]interface{}{
		"name":    c.name,
		"version": c.version,
	})

	if !isUpgrade {
		_, err := c.manager.Define(c.name, c.version, c.location)
		if err != nil {
			return errors.Wrapf(err, "failed to define command %s", c.name)
		}
	}

	err := c.manager.Activate(c.name, c.version)
	if err != nil {
		return errors.Wrapf(err, "failed to activate command %s", c.name)
	}

	return c.removeLegacies(safeVersion)
}

func NewCmdrUpdater(manager core.CommandManager, name, version, localtion string) *CmdrUpdater {
	return &CmdrUpdater{
		manager:  manager,
		name:     name,
		version:  version,
		location: localtion,
	}
}

func init() {
	core.RegisterInitializerFactory("cmdr-updater", func(cfg core.Configuration) (core.Initializer, error) {
		location, err := os.Executable()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get executable location")
		}

		manager, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create command manager")
		}

		return NewCmdrUpdater(manager, core.Name, core.Version, location), nil
	})
}
