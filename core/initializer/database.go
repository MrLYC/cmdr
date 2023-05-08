package initializer

import (
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

type DatabaseMigrator struct {
	dbFactory func() (core.Database, error)
	models    map[core.ModelType]interface{}
}

func (m *DatabaseMigrator) Init(isUpgrade bool) error {
	logger := core.GetLogger()

	db, err := m.dbFactory()
	if err != nil {
		return errors.Wrapf(err, "open database failed")
	}
	defer utils.CallClose(db)

	for name, model := range m.models {
		logger.Debug("initializing database model", map[string]interface{}{
			"model": name,
		})
		err := db.Init(model)
		if err != nil {
			return errors.Wrapf(err, "init database failed")
		}

		logger.Debug("indexing database model", map[string]interface{}{
			"model": name,
		})
		err = db.ReIndex(model)
		if err != nil {
			return errors.Wrapf(err, "reindex database failed")
		}
	}

	return nil
}

func NewDatabaseMigrator(dbFactory func() (core.Database, error), models map[core.ModelType]interface{}) *DatabaseMigrator {
	return &DatabaseMigrator{
		dbFactory: dbFactory,
		models:    models,
	}
}

func init() {
	core.RegisterInitializerFactory("database-migrator", func(cfg core.Configuration) (core.Initializer, error) {
		return NewDatabaseMigrator(func() (core.Database, error) {
			return core.GetDatabase()
		}, core.GetDatabaseModels()), nil
	})
}
