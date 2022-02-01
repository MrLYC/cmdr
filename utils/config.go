package utils

import (
	"fmt"
	"path/filepath"
)

func GetCommandShimsDir(shimsDir, name string) string {
	return filepath.Join(shimsDir, name)
}

func GetCommandShimsPath(shimsDir, name, version string) string {
	return filepath.Join(GetCommandShimsDir(shimsDir, name), fmt.Sprintf("%s_%s", name, version))
}

func GetCommandBinPath(binDir, name string) string {
	return filepath.Join(binDir, name)
}
