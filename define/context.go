package define

//go:generate stringer -type=ContextKey
type ContextKey int

const (
	ContextKeyName ContextKey = iota
	ContextKeyVersion
	ContextKeyLocation
	ContextKeyDBClient
	ContextKeyCommandManaged
	ContextKeyCommands
	ContextKeyCommand
)
