package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewDefineRunner(cfg define.Configuration) Runner {
	return New(
		operator.NewDBClientMaker(),
		operator.NewCommandDefiner(
			cfg.GetString(config.CfgKeyCommandInstallName),
			cfg.GetString(config.CfgKeyCommandInstallVersion),
			cfg.GetString(config.CfgKeyCommandInstallLocation),
			false,
		),
	)
}
