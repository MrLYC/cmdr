package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewUseRunner(cfg define.Configuration) Runner {
	return New(
		operator.NewDBClientMaker(),
		operator.NewSimpleCommandsQuerier(
			cfg.GetString(config.CfgKeyCommandUseName),
			cfg.GetString(config.CfgKeyCommandUseVersion),
		).StrictMode(),
		operator.NewCommandDeactivator(),
		operator.NewBinariesActivator(),
		operator.NewCommandActivator(),
	)
}
