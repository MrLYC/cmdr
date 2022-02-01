package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewUnsetRunner(cfg define.Configuration) Runner {
	return New(
		operator.NewDBClientMaker(),
		operator.NewNamedCommandsQuerier(cfg.GetString(config.CfgKeyCommandUnsetName)),
		operator.NewCommandDeactivator(),
	)
}
