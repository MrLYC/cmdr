package core

import (
	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
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
