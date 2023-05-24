package initializer

import (
	"os"

	ver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type CmdrUpdater struct {
	name     string
	version  string
	location string
	manager  core.CommandManager
}

func (c *CmdrUpdater) collectLegacyVersions() ([]string, error) {
	logger := core.GetLogger()
	currentVersion := ver.Must(ver.NewVersion(c.version))

	query, err := c.manager.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create command query")
	}

	commands, err := query.
		WithName(c.name).
		All()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get commands")
	}

	legacyVersions := make([]string, 0, len(commands))
	for _, command := range commands {
		logger.Debug("checking command", map[string]interface{}{
			"command": command,
		})

		if command.GetActivated() {
			continue
		}

		definedVersion := command.GetVersion()
		semver := ver.Must(ver.NewVersion(definedVersion))

		if currentVersion.Compare(semver) <= 0 {
			continue
		}

		logger.Info("collected legacy cmdr", map[string]interface{}{
			"command": command,
		})
		legacyVersions = append(legacyVersions, definedVersion)
	}

	return legacyVersions, nil
}

func (c *CmdrUpdater) Init(isUpgrade bool) error {
	logger := core.GetLogger()
	logger.Debug("update command", map[string]interface{}{
		"name":    c.name,
		"version": c.version,
	})

	legacyVersions, err := c.collectLegacyVersions()
	if err != nil {
		return errors.Wrapf(err, "failed to collect legacy versions")
	}

	if !isUpgrade {
		_, err := c.manager.Define(c.name, c.version, c.location)
		if err != nil {
			return errors.Wrapf(err, "failed to define command %s", c.name)
		}
	}

	err = c.manager.Activate(c.name, c.version)
	if err != nil {
		return errors.Wrapf(err, "failed to activate command %s", c.name)
	}

	for _, version := range legacyVersions {
		err = c.manager.Undefine(c.name, version)
		if err != nil {
			return errors.Wrapf(err, "failed to undefine command %s", c.name)
		}
	}

	return nil
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
