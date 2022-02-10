package cmdr

type Initializer interface {
	Init() error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Initializer
