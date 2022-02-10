package cmdr

//go:generate stringer -type=CommandProvider

type CommandProvider int

const (
	CommandProviderUnknown CommandProvider = iota
	CommandProviderDatabase
	CommandProviderBinary
	CommandProviderSimple
)

type Command interface {
	Name() string
	Version() string
	Activated() bool
	Location() string
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

	Define(name string, version string, location string) error
	Undefine(name string, version string) error
	Activate(name string, version string) error
	Deactivate(name string) error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Command,CommandQuery,CommandManager
