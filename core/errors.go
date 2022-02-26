package core

import "fmt"

var (
	ErrCommandAlreadyActivated = fmt.Errorf("command already activated")
	ErrShellNotSupported       = fmt.Errorf("shell not supported")
	ErrBinaryNotFound          = fmt.Errorf("binaries not found")
	ErrReleaseAssetNotFound    = fmt.Errorf("release asset not found")
)
