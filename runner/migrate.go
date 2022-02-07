package runner

import (
	"github.com/mrlyc/cmdr/core/model"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

func NewMigrateRunner(cfg define.Configuration, helper *utils.CmdrHelper) *OperatorRunner {
	return New(
		operator.NewDirectoryMaker(map[string]string{
			"shims": helper.GetShimsDir(),
			"bin":   helper.GetBinDir(),
		}),
		operator.NewDBClientMaker(helper),
		operator.NewDBMigrator(new(model.Command)),
	)
}
