package core

import "fmt"

//go:generate stringer -type=CommandProvider

type CommandProvider int

const (
	CommandProviderUnknown CommandProvider = iota
	CommandProviderDefault
	CommandProviderDatabase
	CommandProviderBinary
	CommandProviderDownload
)

type Command interface {
	GetName() string
	GetVersion() string
	GetActivated() bool
	GetLocation() string
}

type CommandQuery interface {
	WithName(name string) CommandQuery
	WithVersion(version string) CommandQuery
	WithActivated(activated bool) CommandQuery
	WithLocation(location string) CommandQuery

	All() ([]Command, error)
	One() (Command, error)
	Count() (int, error)
}

type CommandManager interface {
	Close() error

	Provider() CommandProvider

	Query() (CommandQuery, error)

	Define(name string, version string, location string) (Command, error)
	Undefine(name string, version string) error
	Activate(name string, version string) error
	Deactivate(name string) error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Command,CommandQuery,CommandManager

var (
	ErrCommandManagerFactoryeNotFound = fmt.Errorf("factory not found")
	factoriesCommandManager           map[CommandProvider]func(cfg Configuration) (CommandManager, error)
)

func GetCommandManagerFactory(key CommandProvider) func(cfg Configuration) (CommandManager, error) {
	fn, ok := factoriesCommandManager[key]

	if !ok {
		return nil
	}

	return fn
}

func RegisterCommandManagerFactory(key CommandProvider, fn func(cfg Configuration) (CommandManager, error)) {
	factoriesCommandManager[key] = fn
}

func NewCommandManager(key CommandProvider, cfg Configuration) (CommandManager, error) {
	fn, ok := factoriesCommandManager[key]

	if !ok {
		return nil, ErrCommandManagerFactoryeNotFound
	}

	return fn(cfg)
}

func init() {
	factoriesCommandManager = make(map[CommandProvider]func(cfg Configuration) (CommandManager, error))
}
