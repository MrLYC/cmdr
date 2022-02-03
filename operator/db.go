package operator

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type StormClient struct {
	*storm.DB
}

func NewDBClient(path string) (define.DBClient, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "open database failed")
	}

	return &StormClient{
		DB: db,
	}, nil
}

type DBClientMaker struct {
	client define.DBClient
	helper *utils.CmdrHelper
}

func (m *DBClientMaker) String() string {
	return "db-client"
}

func (m *DBClientMaker) Run(ctx context.Context) (context.Context, error) {
	path := m.helper.GetDatabasePath()
	client, err := NewDBClient(path)
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

func NewDBClientMaker(helper *utils.CmdrHelper) *DBClientMaker {
	return &DBClientMaker{
		helper: helper,
	}
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
