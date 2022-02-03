package config

import "github.com/mrlyc/cmdr/define"

const (
	// command.list
	CfgKeyCommandListName     = "command.list.name"
	CfgKeyCommandListVersion  = "command.list.version"
	CfgKeyCommandListLocation = "command.list.location"
	CfgKeyCommandListActivate = "command.list.activate"
)

func GetCommandListName(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandListName)
}

func GetCommandListVersion(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandListVersion)
}

func GetCommandListLocation(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCommandListLocation)
}

func GetCommandListActivate(cfg define.Configuration) bool {
	return cfg.GetBool(CfgKeyCommandListActivate)
}
