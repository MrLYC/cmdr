package config

import "github.com/mrlyc/cmdr/define"

const (
	// command.install
	CfgKeyCommandInstallName     = "command.install.name"
	CfgKeyCommandInstallVersion  = "command.install.version"
	CfgKeyCommandInstallLocation = "command.install.location"
	CfgKeyCommandInstallActivate = "command.install.activate"
)

func GetCommandInstallName(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandInstallName)
}

func GetCommandInstallVersion(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandInstallVersion)
}

func GetCommandInstallLocation(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandInstallLocation)
}

func GetCommandInstallActivate(cfg define.Configuration) bool {
	return cfg.GetBool(CfgKeyCommandInstallActivate)
}
