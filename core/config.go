package core

import (
	"fmt"
	"path/filepath"

	"github.com/mrlyc/cmdr/define"
)

func GetRootDir() string {
	return define.Configuration.GetString("cmdr.root")
}

func GetBinDir() string {
	return filepath.Join(GetRootDir(), "bin")
}

func GetShimsDir() string {
	return filepath.Join(GetRootDir(), "shims")
}

func GetDBName() string {
	return define.Configuration.GetString("database.name")
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
