package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewDefineRunner(cfg define.Configuration) Runner {
	shimsDir := config.GetShimsDir()

	return New(
		operator.NewDBClientMaker(),
		operator.NewCommandDefiner(
			shimsDir,
			cfg.GetString(config.CfgKeyCommandInstallName),
			cfg.GetString(config.CfgKeyCommandInstallVersion),
			cfg.GetString(config.CfgKeyCommandInstallLocation),
			false,
		),
	)
}
