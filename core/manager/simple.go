package manager

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type SimpleManager struct {
	main     core.CommandManager
	recorder core.CommandManager
}

func (m *SimpleManager) each(fn func(mgr core.CommandManager) error) error {
	for _, mgr := range []core.CommandManager{m.main, m.recorder} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *SimpleManager) reverseEach(fn func(mgr core.CommandManager) error) error {
	for _, mgr := range []core.CommandManager{m.recorder, m.main} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *SimpleManager) all(fn func(mgr core.CommandManager) error) error {
	var errs error
	for _, mgr := range []core.CommandManager{m.main, m.recorder} {
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
	return m.recorder.Query()
}

func (m *SimpleManager) Define(name, version, location string) (core.Command, error) {
	binary, err := m.main.Define(name, version, location)
	if err != nil {
		return nil, err
	}

	return m.recorder.Define(name, version, binary.GetLocation())
}

func (m *SimpleManager) Undefine(name, version string) error {
	return m.reverseEach(func(mgr core.CommandManager) error {
		return mgr.Undefine(name, version)
	})
}

func (m *SimpleManager) Activate(name, version string) error {
	return m.each(func(mgr core.CommandManager) error {
		return mgr.Activate(name, version)
	})
}

func (m *SimpleManager) Deactivate(name string) error {
	return m.reverseEach(func(mgr core.CommandManager) error {
		return mgr.Deactivate(name)
	})
}

func NewSimpleManager(main core.CommandManager, recorder core.CommandManager) *SimpleManager {
	return &SimpleManager{main: main, recorder: recorder}
}

func init() {
	var _ core.CommandManager = (*SimpleManager)(nil)

	core.RegisterCommandManagerFactory(core.CommandProviderDefault, func(cfg core.Configuration) (core.CommandManager, error) {
		mainMgr, err := core.NewCommandManager(core.CommandProviderBinary, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "new main command manager failed")
		}

		recorderMgr, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "new recorder command manager failed")
		}

		return NewSimpleManager(mainMgr, recorderMgr), nil
	})
}
