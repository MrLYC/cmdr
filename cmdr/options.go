package cmdr

type Option interface {
	Apply(target Optional) error
}

type Optional interface {
	ApplyOptions(options ...Option) error
}
