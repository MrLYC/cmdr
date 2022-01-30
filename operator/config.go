package operator

import (
	"fmt"
	"path/filepath"

	"github.com/mrlyc/cmdr/define"
)

func GetRootDir() string {
	return define.Config.GetString(define.CfgKeyCmdrRoot)
}

func GetBinDir() string {
	return filepath.Join(GetRootDir(), define.Config.GetString(define.CfgKeyBinDir))
}

func GetShimsDir() string {
	return filepath.Join(GetRootDir(), define.Config.GetString(define.CfgKeyShimsDir))
}

func GetDatabasePath() string {
	return filepath.Join(GetRootDir(), define.Config.GetString(define.CfgKeyDatabase))
}

func GetCommandShimsDir(shimsDir, name string) string {
	return filepath.Join(shimsDir, name)
}

func GetCommandShimsPath(shimsDir, name, version string) string {
	return filepath.Join(GetCommandShimsDir(shimsDir, name), fmt.Sprintf("%s_%s", name, version))
}

func GetCommandBinPath(binDir, name string) string {
	return filepath.Join(binDir, name)
}
