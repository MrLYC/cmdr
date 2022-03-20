package core

import "github.com/asdine/storm/v3"

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

func init() {
	databaseModels = make(map[ModelType]interface{})
}
