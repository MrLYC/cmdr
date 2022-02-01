package operator

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
)

//go:generate mockgen -destination=mock/storm.go -package=mock github.com/asdine/storm/v3 Query

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock DBClient
type DBClient interface {
	storm.TypeStore
	Close() error
}

type StormClient struct {
	*storm.DB
}

func NewDBClient(path string) (DBClient, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "open database failed")
	}

	return &StormClient{
		DB: db,
	}, nil
}

func GetDBClient() (DBClient, error) {
	db := config.GetDatabasePath()
	logger := define.Logger

	logger.Debug("opening database", map[string]interface{}{
		"name": db,
	})

	return NewDBClient(db)
}

type DBClientMaker struct {
	client DBClient
}

func (m *DBClientMaker) String() string {
	return "db-client"
}

func (m *DBClientMaker) Run(ctx context.Context) (context.Context, error) {
	client, err := GetDBClient()
	if err != nil {
		return ctx, errors.Wrapf(err, "create database client failed")
	}

	m.client = client

	return context.WithValue(ctx, define.ContextKeyDBClient, m.client), nil
}

func (m *DBClientMaker) Commit(ctx context.Context) error {
	return m.client.Close()
}

func (m *DBClientMaker) Rollback(ctx context.Context) {
	_ = m.client.Close()
}

func NewDBClientMaker() *DBClientMaker {
	return &DBClientMaker{}
}

type DBMigrator struct {
	BaseOperator
	models []interface{}
}

func (m *DBMigrator) String() string {
	return "db-migrator"
}

func (m *DBMigrator) Run(ctx context.Context) (context.Context, error) {
	client := GetDBClientFromContext(ctx)
	logger := define.Logger

	logger.Debug("database migrating")
	for _, model := range m.models {
		logger.Debug("migrating model", map[string]interface{}{
			"model": model,
		})
		err := client.Init(model)
		if err != nil {
			return ctx, errors.Wrapf(err, "migrate model %T failed", model)
		}

		err = client.ReIndex(model)
		if err != nil {
			return ctx, errors.Wrapf(err, "indexing model %T failed", model)
		}
	}

	return ctx, nil
}

func NewDBMigrator(models ...interface{}) *DBMigrator {
	return &DBMigrator{
		models: models,
	}
}
