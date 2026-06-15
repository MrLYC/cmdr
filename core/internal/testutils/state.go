package testutils

import "github.com/mrlyc/cmdr/core"

func SwapConfiguration(cfg core.Configuration) func() {
	previous := core.GetConfiguration()
	core.SetConfiguration(cfg)
	return func() {
		core.SetConfiguration(previous)
	}
}

func RestoreCommandManagerFactory(provider core.CommandProvider) func() {
	previous := core.GetCommandManagerFactory(provider)
	return func() {
		core.RegisterCommandManagerFactory(provider, previous)
	}
}

func RegisterCommandManagerFactory(provider core.CommandProvider, factory func(core.Configuration) (core.CommandManager, error)) func() {
	restore := RestoreCommandManagerFactory(provider)
	core.RegisterCommandManagerFactory(provider, factory)
	return restore
}
