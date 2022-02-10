package core

import "fmt"

type factoryCommandManager func(cfg Configuration, opts ...Option) (CommandManager, error)

var (
	ErrCommandManagerFactoryeNotFound = fmt.Errorf("factory not found")
	factoriesCommandManager           map[CommandProvider]factoryCommandManager
)

func RegisterCommandManagerFactory(key CommandProvider, fn factoryCommandManager) {
	factoriesCommandManager[key] = fn
}

func NewCommandManager(key CommandProvider, cfg Configuration, opts ...Option) (CommandManager, error) {
	fn, ok := factoriesCommandManager[key]

	if !ok {
		return nil, ErrCommandManagerFactoryeNotFound
	}

	return fn(cfg, opts...)
}

func init() {
	factoriesCommandManager = make(map[CommandProvider]factoryCommandManager)
}
