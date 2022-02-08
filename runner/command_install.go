package runner

import (
	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/operator"
)

func NewInstallRunner(cfg define.Configuration, cmdr *core.Cmdr) define.Runner {
	runner := New(
		operator.NewCommandDefiner(
			cmdr,
			cfg.GetString(config.CfgKeyCommandInstallName),
			cfg.GetString(config.CfgKeyCommandInstallVersion),
			cfg.GetString(config.CfgKeyCommandInstallLocation),
		),
		operator.NewDownloader(),
		operator.NewBinariesChecker(),
		operator.NewBinariesInstaller(cmdr, true),
	)

	if cfg.GetBool(config.CfgKeyCommandInstallActivate) {
		runner.Add(
			operator.NewCommandDeactivator(cmdr),
			operator.NewBinariesActivator(cmdr),
			operator.NewCommandActivator(cmdr),
		)
	}

	return runner
}
