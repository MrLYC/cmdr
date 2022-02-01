package config

import (
	"path/filepath"

	"github.com/mrlyc/cmdr/define"
)

func GetCmdrRoot(cfg define.Configuration) string {
	return cfg.GetString(CfgKeyCmdrRoot)
}

func GetBinDir(cfg define.Configuration) string {
	return filepath.Join(GetCmdrRoot(cfg), Global.GetString(CfgKeyBinDir))
}

func GetShimsDir(cfg define.Configuration) string {
	return filepath.Join(GetCmdrRoot(cfg), Global.GetString(CfgKeyShimsDir))
}

func GetDatabasePath(cfg define.Configuration) string {
	return filepath.Join(GetCmdrRoot(cfg), Global.GetString(CfgKeyDatabase))
}
