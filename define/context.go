package define

//go:generate stringer -type=ContextKey
type ContextKey int

const (
	ContextKeyDBClient ContextKey = iota
	ContextKeyCommands
)
