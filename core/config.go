package core

import (
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
