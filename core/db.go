package core

import (
	"path"

	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type Client = *StormClient

type StormClient struct {
	*storm.DB
}

func (c *StormClient) Atomic(fn func() error) error {
	return fn()
}

func (c *StormClient) Migrate(models ...interface{}) (err error) {
	for _, model := range models {
		err = c.Init(model)
		if err != nil {
			return errors.Wrapf(err, "migrate model %T failed", model)
		}

		err = c.ReIndex(model)
		if err != nil {
			return errors.Wrapf(err, "indexing model %T failed", model)
		}
	}

	return nil
}

func GetClient() Client {
	cfg := define.Configuration
	cmdrDir := cfg.GetString("cmdr.root")
	name := cfg.GetString("database.name")

	db, err := storm.Open(path.Join(cmdrDir, name))
	utils.CheckError(err)
	return &StormClient{
		DB: db,
	}
}

func IsQueryNotFound(err error) bool {
	return errors.Cause(err) == storm.ErrNotFound
}
