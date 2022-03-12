package utils

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

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

type CobraCommandCompleteHelper struct {
	managerProvider core.CommandProvider
	commands        []core.Command
	queryOnce       sync.Once
	cobraCommand    *cobra.Command
	flagName        string
	flagVersion     string
	flagLocation    string
	flagActivate    string
}

func (h *CobraCommandCompleteHelper) updateQuery(query core.CommandQuery) core.CommandQuery {
	flags := h.cobraCommand.Flags()
	name, err := flags.GetString(h.flagName)
	if err == nil && name != "" {
		query.WithName(name)
	}

	version, err := flags.GetString(h.flagVersion)
	if err == nil && version != "" {
		query.WithVersion(version)
	}

	location, err := flags.GetString(h.flagLocation)
	if err == nil && location != "" {
		query.WithLocation(location)
	}

	activate, err := flags.GetBool(h.flagActivate)
	if err == nil && activate {
		query.WithActivated(activate)
	}

	return query
}

func (h *CobraCommandCompleteHelper) getCommands() []core.Command {
	logger := core.GetLogger()

	h.queryOnce.Do(func() {
		manager, err := core.NewCommandManager(h.managerProvider, core.GetConfiguration())
		if err != nil {
			logger.Debug("Failed to create command manager", map[string]interface{}{
				"error": err,
			})
			return
		}

		defer manager.Close()

		query, err := manager.Query()
		if err != nil {
			logger.Debug("Failed to create command query", map[string]interface{}{
				"error": err,
			})
			return
		}

		h.commands, err = h.updateQuery(query).All()
		if err != nil {
			logger.Debug("Failed to query commands", map[string]interface{}{
				"error": err,
			})
			return
		}
	})

	return h.commands
}

func (h *CobraCommandCompleteHelper) isFlagSet(name string) bool {
	flags := h.cobraCommand.Flags()
	return flags.Lookup(name) != nil
}

func (h *CobraCommandCompleteHelper) GetNameSlice(prefix string) []string {
	commands := h.getCommands()
	results := make([]string, 0, len(commands))

	for _, command := range commands {
		name := command.GetName()
		if strings.HasPrefix(name, prefix) {
			results = append(results, name)
		}
	}

	return results
}

func (h *CobraCommandCompleteHelper) GetVersionSlice(prefix string) []string {
	commands := h.getCommands()
	results := make([]string, 0, len(commands))

	for _, command := range commands {
		version := command.GetVersion()
		if strings.HasPrefix(version, prefix) {
			results = append(results, version)
		}
	}

	return results
}

func (h *CobraCommandCompleteHelper) GetLocationSlice(prefix string) []string {
	commands := h.getCommands()
	results := make([]string, 0, len(commands))

	for _, command := range commands {
		location := command.GetLocation()
		if strings.HasPrefix(location, prefix) {
			results = append(results, location)
		}
	}

	return results
}

func (h *CobraCommandCompleteHelper) RegisterNameFunc() error {
	return h.cobraCommand.RegisterFlagCompletionFunc(
		h.flagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return h.GetNameSlice(toComplete), cobra.ShellCompDirectiveDefault
		},
	)
}

func (h *CobraCommandCompleteHelper) RegisterVersionFunc() error {
	return h.cobraCommand.RegisterFlagCompletionFunc(
		h.flagVersion,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return h.GetVersionSlice(toComplete), cobra.ShellCompDirectiveDefault
		},
	)
}

func (h *CobraCommandCompleteHelper) RegisterLocationFunc() error {
	return h.cobraCommand.RegisterFlagCompletionFunc(
		h.flagLocation,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return h.GetLocationSlice(toComplete), cobra.ShellCompDirectiveDefault
		},
	)
}

func (h *CobraCommandCompleteHelper) RegisterAll() error {
	mappings := map[string]func() error{
		h.flagName:     h.RegisterNameFunc,
		h.flagVersion:  h.RegisterVersionFunc,
		h.flagLocation: h.RegisterLocationFunc,
	}

	for name, fn := range mappings {
		if h.isFlagSet(name) {
			err := fn()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewCobraCommandCompleteHelper(cmd *cobra.Command, provider core.CommandProvider) *CobraCommandCompleteHelper {
	return &CobraCommandCompleteHelper{
		managerProvider: provider,
		cobraCommand:    cmd,
		commands:        make([]core.Command, 0),
		flagName:        "name",
		flagVersion:     "version",
		flagLocation:    "location",
		flagActivate:    "activate",
	}
}

func NewDefaultCobraCommandCompleteHelper(cmd *cobra.Command) *CobraCommandCompleteHelper {
	return NewCobraCommandCompleteHelper(cmd, core.CommandProviderDefault)
}
