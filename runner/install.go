package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

func NewInstallRunner(cfg define.Configuration, helper *utils.CmdrHelper) Runner {
	runner := New(
		operator.NewDBClientMaker(helper),
		operator.NewCommandDefiner(
			cfg.GetString(config.CfgKeyCommandInstallName),
			cfg.GetString(config.CfgKeyCommandInstallVersion),
			cfg.GetString(config.CfgKeyCommandInstallLocation),
			true,
			helper,
		),
		operator.NewDownloader(),
		operator.NewBinariesInstaller(helper),
	)

	if cfg.GetBool(config.CfgKeyCommandInstallActivate) {
		runner.Add(
			operator.NewCommandDeactivator(),
			operator.NewBinariesActivator(helper),
			operator.NewCommandActivator(),
		)
	}

	return runner
}
