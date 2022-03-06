package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

var ErrCmdrCommandAlreadyDefined = errors.New("cmdr command already defined")

func normalizeVersion(version string) string {
	return strings.TrimPrefix(version, "v")
}

func DefineCmdrCommand(manager core.CommandManager, name string, version string, location string, activate bool) (core.Command, error) {
	version = normalizeVersion(version)
	command, err := manager.Define(name, version, location)
	if err != nil {
		return nil, err
	}

	if activate {
		return command, manager.Activate(name, version)
	}

	return command, nil
}

func DefineCmdrCommandNX(manager core.CommandManager, name string, version string, location string, activate bool) (core.Command, error) {
	query, err := manager.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "query command %v failed", name)
	}

	_, err = query.WithName(name).WithVersion(version).One()
	if err == nil {
		return nil, errors.Wrapf(ErrCmdrCommandAlreadyDefined, "%v(%s)", name, version)
	}

	return DefineCmdrCommand(manager, name, version, location, activate)
}

func GetCmdrCommand(manager core.CommandManager, name, version string) (core.Command, error) {
	query, err := manager.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "query command %v failed", name)
	}

	command, err := query.WithName(name).WithVersion(version).One()
	if err != nil {
		return nil, errors.Wrapf(err, "query command %v(%s) failed", name, version)
	}

	return command, nil
}

func UpgradeCmdr(ctx context.Context, cfg core.Configuration, url, version string, args []string) error {
	name := core.Name
	version = normalizeVersion(version)
	manager, err := core.NewCommandManager(core.CommandProviderDownload, cfg)
	if err != nil {
		return errors.Wrapf(err, "create command manager %v failed", core.CommandProviderDownload)
	}

	_, err = DefineCmdrCommandNX(manager, name, version, url, false)
	if err != nil {
		return errors.Wrapf(err, "define command %v failed", name)
	}

	command, err := GetCmdrCommand(manager, name, version)
	if err != nil {
		return errors.Wrapf(err, "get command %v failed", name)
	}

	err = manager.Close()
	if err != nil {
		return errors.Wrapf(err, "close command manager %v failed", core.CommandProviderDownload)
	}

	err = WaitProcess(ctx, command.GetLocation(), args)
	if err != nil {
		return errors.Wrapf(err, "run command %v failed", name)
	}

	return nil
}

func RunCobraCommandWith(provider core.CommandProvider, fn func(cfg core.Configuration, manager core.CommandManager) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg := core.GetConfiguration()

		manager, err := core.NewCommandManager(provider, cfg)
		if err != nil {
			ExitOnError("Failed to create command manager", err)
		}

		defer CallClose(manager)

		ExitOnError(fmt.Sprintf("Failed to run command %s", cmd.Name()), fn(cfg, manager))
	}
}
