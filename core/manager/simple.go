package manager

import (
	"github.com/hashicorp/go-multierror"
	ver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type SimpleManager struct {
	main      core.CommandManager
	followers []core.CommandManager
}

func (m *SimpleManager) each(fn func(mgr core.CommandManager) error) error {
	err := fn(m.main)
	if err != nil {
		return err
	}

	for _, mgr := range m.followers {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *SimpleManager) all(fn func(mgr core.CommandManager) error) error {
	var errs error

	err := fn(m.main)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	for _, mgr := range m.followers {
		err := fn(mgr)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (m *SimpleManager) Close() error {
	return m.all(func(mgr core.CommandManager) error {
		return mgr.Close()
	})
}

func (m *SimpleManager) Provider() core.CommandProvider {
	return core.CommandProviderDefault
}

func (m *SimpleManager) Query() (core.CommandQuery, error) {
	return m.main.Query()
}

func (m *SimpleManager) normalizeVersion(version string) (string, error) {
	if version == "" {
		return "", nil
	}

	v, err := ver.NewVersion(version)
	if err != nil {
		return "", errors.Wrapf(err, "invalid version %s", version)
	}

	return v.String(), nil
}

func (m *SimpleManager) Define(name, version, location string) (core.Command, error) {
	semver, err := m.normalizeVersion(version)
	if err != nil {
		return nil, err
	}

	var result core.Command
	return result, m.each(func(mgr core.CommandManager) error {
		command, err := mgr.Define(name, semver, location)
		if command == nil {
			result = command
		}

		return err
	})
}

func (m *SimpleManager) Undefine(name, version string) error {
	return m.each(func(mgr core.CommandManager) error {
		return mgr.Undefine(name, version)
	})
}

func (m *SimpleManager) Activate(name, version string) error {
	return m.each(func(mgr core.CommandManager) error {
		return mgr.Activate(name, version)
	})
}

func (m *SimpleManager) Deactivate(name string) error {
	return m.each(func(mgr core.CommandManager) error {
		return mgr.Deactivate(name)
	})
}

func NewSimpleManager(main core.CommandManager, followers []core.CommandManager) *SimpleManager {
	return &SimpleManager{main: main, followers: followers}
}

func init() {
	var _ core.CommandManager = (*SimpleManager)(nil)

	core.RegisterCommandManagerFactory(core.CommandProviderDefault, func(cfg core.Configuration) (core.CommandManager, error) {
		mainMgr, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "new main command manager failed")
		}

		followers := make([]core.CommandManager, 0, 1)
		for _, provider := range []core.CommandProvider{} {
			mgr, err := core.NewCommandManager(provider, cfg)
			if err != nil {
				return nil, errors.Wrapf(err, "new manager %v failed", provider)
			}

			followers = append(followers, mgr)
		}

		return NewSimpleManager(mainMgr, followers), nil
	})
}
