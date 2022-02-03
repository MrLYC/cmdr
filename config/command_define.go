package config

import "github.com/mrlyc/cmdr/define"

const (
	// command.define
	CfgKeyCommandDefineName     = "command.define.name"
	CfgKeyCommandDefineVersion  = "command.define.version"
	CfgKeyCommandDefineLocation = "command.define.location"
)

func GetCommandDefineName(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandDefineName)
}

func GetCommandDefineVersion(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandDefineVersion)
}

func GetCommandDefineLocation(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandDefineLocation)
}
