package manager

import (
	"github.com/hashicorp/go-multierror"

	"github.com/mrlyc/cmdr/cmdr"
)

type SimpleManager struct {
	main    cmdr.CommandManager
	recoder cmdr.CommandManager
}

func (m *SimpleManager) each(fn func(mgr cmdr.CommandManager) error) error {
	for _, mgr := range []cmdr.CommandManager{m.main, m.recoder} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *SimpleManager) reverseEach(fn func(mgr cmdr.CommandManager) error) error {
	for _, mgr := range []cmdr.CommandManager{m.recoder, m.main} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *SimpleManager) Init() error {
	var errs error
	for _, mgr := range []cmdr.CommandManager{m.main, m.recoder} {
		err := mgr.Init()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (m *SimpleManager) Provider() cmdr.CommandProvider {
	return cmdr.CommandProviderSimple
}

func (m *SimpleManager) Query() (cmdr.CommandQuery, error) {
	return m.recoder.Query()
}

func (m *SimpleManager) Define(name, version, location string) error {
	return m.each(func(mgr cmdr.CommandManager) error {
		return mgr.Define(name, version, location)
	})
}

func (m *SimpleManager) Undefine(name, version string) error {
	return m.reverseEach(func(mgr cmdr.CommandManager) error {
		return mgr.Undefine(name, version)
	})
}

func (m *SimpleManager) Activate(name, version string) error {
	return m.each(func(mgr cmdr.CommandManager) error {
		return mgr.Activate(name, version)
	})
}

func (m *SimpleManager) Deactivate(name string) error {
	return m.reverseEach(func(mgr cmdr.CommandManager) error {
		return mgr.Deactivate(name)
	})
}

func NewSimpleManager(main cmdr.CommandManager, shadow cmdr.CommandManager) *SimpleManager {
	return &SimpleManager{main: main, recoder: shadow}
}

func init() {
	var _ cmdr.CommandManager = (*SimpleManager)(nil)
}
