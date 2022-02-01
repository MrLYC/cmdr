package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewInstallRunner(cfg define.Configuration) Runner {
	runner := New(
		operator.NewDBClientMaker(),
		operator.NewCommandDefiner(
			cfg.GetString(config.CfgKeyCommandInstallName),
			cfg.GetString(config.CfgKeyCommandInstallVersion),
			cfg.GetString(config.CfgKeyCommandInstallLocation),
			true,
		),
		operator.NewDownloader(),
		operator.NewBinariesInstaller(),
	)

	if cfg.GetBool(config.CfgKeyCommandInstallActivate) {
		runner.Add(
			operator.NewCommandDeactivator(),
			operator.NewBinariesActivator(),
			operator.NewCommandActivator(),
		)
	}

	return runner
}
