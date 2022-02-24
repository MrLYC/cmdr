package initializer

import (
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type CmdrUpdater struct {
	manager core.CommandManager
}

func (c *CmdrUpdater) getActivatedCmdrVersion() string {
	query, err := c.manager.Query()
	if err != nil {
		return ""
	}

	command, err := query.
		WithName(core.Name).
		WithActivated(true).
		One()

	if err != nil {
		return ""
	}

	return command.GetVersion()
}

func (c *CmdrUpdater) install() error {
	location, err := os.Executable()
	if err != nil {
		return errors.Wrapf(err, "failed to get executable location")
	}

	err = c.manager.Define(core.Name, core.Version, location)
	if err != nil {
		return errors.Wrapf(err, "failed to define command %s", core.Name)
	}

	err = c.manager.Activate(core.Name, core.Version)
	if err != nil {
		return errors.Wrapf(err, "failed to activate command %s", core.Name)
	}

	return nil
}

func (c *CmdrUpdater) removeLegacies(safeVersions []string) error {
	query, err := c.manager.Query()
	if err != nil {
		return errors.Wrapf(err, "failed to create command query")
	}

	commands, err := query.
		WithName(core.Name).
		All()

	if err != nil {
		return errors.Wrapf(err, "failed to get commands")
	}

	var errs error
	for _, command := range commands {
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

		err = c.manager.Undefine(core.Name, version)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (c *CmdrUpdater) Init() error {
	safeVersion := []string{core.Version}
	version := c.getActivatedCmdrVersion()
	if version != "" {
		safeVersion = append(safeVersion, version)
	}

	err := c.install()
	if err != nil {
		return errors.WithMessagef(err, "failed to install command %s", core.Name)
	}

	return c.removeLegacies(safeVersion)
}

func NewCmdrUpdater(manager core.CommandManager) *CmdrUpdater {
	return &CmdrUpdater{
		manager: manager,
	}
}

func init() {
	core.RegisterInitializerFactory("cmdr-updater", func(cfg core.Configuration) (core.Initializer, error) {
		manager, err := core.NewCommandManager(core.CommandProviderDefault, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create command manager")
		}

		return NewCmdrUpdater(manager), nil
	})
}
