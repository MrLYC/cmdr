package config

import (
	"path/filepath"
)

func GetCmdrRoot() string {
	return Global.GetString(CfgKeyCmdrRoot)
}

func GetBinDir() string {
	return filepath.Join(GetCmdrRoot(), Global.GetString(CfgKeyBinDir))
}

func GetShimsDir() string {
	return filepath.Join(GetCmdrRoot(), Global.GetString(CfgKeyShimsDir))
}

func GetDatabasePath() string {
	return filepath.Join(GetCmdrRoot(), Global.GetString(CfgKeyDatabase))
}
