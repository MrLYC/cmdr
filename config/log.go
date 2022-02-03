package config

import "github.com/mrlyc/cmdr/define"

const (
	// log
	CfgKeyLogLevel  = "log.level"
	CfgKeyLogOutput = "log.output"
)

func GetLogLevel(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyLogLevel)
}

func GetLogOutput(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyLogOutput)
}
