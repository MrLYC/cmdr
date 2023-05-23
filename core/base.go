package core

import (
	"fmt"
	"runtime"
)

var (
	Author    = "MrLYC"
	Name      = "cmdr"
	Version   = "0.0.0"
	Commit    = ""
	BuildDate = ""
	Asset     = ""
)

func init() {
	if Asset == "" {
		Asset = fmt.Sprintf("%s_%s_%s", Name, runtime.GOOS, runtime.GOARCH)
	}
}
