package runner

import (
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

const (
	CfgKeyCommandUnsetName = "command.unset.name"
)

func NewUnsetRunner(cfg define.Configuration) Runner {
	return New(
		operator.NewDBClientMaker(),
		operator.NewNamedCommandsQuerier(cfg.GetString(CfgKeyCommandUnsetName)),
		operator.NewCommandDeactivator(),
	)
}
