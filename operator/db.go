package operator

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type DBMigrator struct {
	*CmdrOperator
}

func (m *DBMigrator) String() string {
	return "db-migrator"
}

func (m *DBMigrator) Run(ctx context.Context) (context.Context, error) {
	err := m.cmdr.Init()
	if err != nil {
		return ctx, errors.Wrapf(err, "initialize cmdr failed")
	}

	return ctx, nil
}

func NewDBMigrator(cmdr *core.Cmdr) *DBMigrator {
	return &DBMigrator{
		CmdrOperator: NewCmdrOperator(cmdr),
	}
}
