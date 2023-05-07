package core

import (
	"errors"

	"github.com/asdine/storm/v3"
)

//go:generate mockgen -destination=mock/storm.go -package=mock github.com/asdine/storm/v3 Query
//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock DBClient

type Database interface {
	storm.TypeStore
	Close() error
}

//go:generate stringer -type=ModelType
type ModelType int

const (
	ModelTypeUnknown ModelType = iota
	ModelTypeCommand
)

var databaseModels map[ModelType]interface{}

func RegisterDatabaseModel(modelType ModelType, model interface{}) {
	databaseModels[modelType] = model
}

func GetDatabaseModels() map[ModelType]interface{} {
	return databaseModels
}

func GetDatabaseModel(modelType ModelType) interface{} {
	return databaseModels[modelType]
}

var (
	databaseFactory          func() (Database, error)
	ErrDatabaseFactoryNotSet = errors.New("database factory not set")
)

func SetDatabaseFactory(fn func() (Database, error)) {
	databaseFactory = fn
}

func GetDatabaseFactory() func() (Database, error) {
	return databaseFactory
}

func GetDatabase() (Database, error) {
	database, err := databaseFactory()
	return database, err
}

func init() {
	databaseModels = make(map[ModelType]interface{})
	databaseFactory = func() (Database, error) {
		return nil, ErrDatabaseFactoryNotSet
	}
}
