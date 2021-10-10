package core

import (
	"path"

	"github.com/mrlyc/cmdr/define"
)

func GetRootDir() string {
	return define.Configuration.GetString("cmdr.root")
}

func GetBinDir() string {
	return path.Join(GetRootDir(), "bin")
}

func GetShimsDir() string {
	return path.Join(GetRootDir(), "shims")
}
