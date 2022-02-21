package core

import "fmt"

type Initializer interface {
	Init() error
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Initializer

type factoryInitializer func(cfg Configuration) (Initializer, error)

var (
	ErrInitializerFactoryeNotFound = fmt.Errorf("factory not found")
	factoriesInitializer           map[string]factoryInitializer
)

func RegisterInitializerFactory(key string, fn func(cfg Configuration) (Initializer, error)) {
	factoriesInitializer[key] = fn
}

func NewInitializer(key string, cfg Configuration) (Initializer, error) {
	fn, ok := factoriesInitializer[key]

	if !ok {
		return nil, ErrInitializerFactoryeNotFound
	}

	return fn(cfg)
}

func init() {
	factoriesInitializer = make(map[string]factoryInitializer)
}
