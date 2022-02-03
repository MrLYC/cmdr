package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

func NewUnsetRunner(cfg define.Configuration, helper *utils.CmdrHelper) Runner {
	return New(
		operator.NewDBClientMaker(helper),
		operator.NewNamedCommandsQuerier(cfg.GetString(config.CfgKeyCommandUnsetName)),
		operator.NewCommandDeactivator(),
	)
}
