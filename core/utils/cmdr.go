package utils

import (
	"strings"

	"github.com/mrlyc/cmdr/core"
)

func DefineCmdrCommand(manager core.CommandManager, name string, version string, location string, activate bool) error {
	version = strings.TrimPrefix(version, "v")
	err := manager.Define(name, version, location)
	if err != nil {
		return err
	}

	if activate {
		return manager.Activate(name, version)
	}

	return nil
}
