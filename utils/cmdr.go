package utils

import (
	"fmt"
	"path/filepath"
)

type CmdrHelper struct {
	root string
}

func (h *CmdrHelper) GetRootDir() string {
	return h.root
}

func (h *CmdrHelper) GetBinDir() string {
	return filepath.Join(h.root, "bin")
}

func (h *CmdrHelper) GetShimsDir() string {
	return filepath.Join(h.root, "shims")
}

func (h *CmdrHelper) GetDatabasePath() string {
	return filepath.Join(h.root, "cmdr.db")
}

func (h *CmdrHelper) GetCommandShimsDir(name string) string {
	return filepath.Join(h.GetShimsDir(), name)
}

func (h *CmdrHelper) GetCommandShimsPath(name, version string) string {
	return filepath.Join(h.GetCommandShimsDir(name), fmt.Sprintf("%s_%s", name, version))
}

func (h *CmdrHelper) GetCommandBinPath(name string) string {
	return filepath.Join(h.GetBinDir(), name)
}

func NewCmdrHelper(root string) *CmdrHelper {
	return &CmdrHelper{root: root}
}
