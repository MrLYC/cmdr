package runner

import (
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewMigrateRunner(cfg define.Configuration, cmdr *core.Cmdr) *OperatorRunner {
	return New(
		operator.NewDBMigrator(cmdr),
	)
}
