package operator

import (
	"fmt"
)

var (
	ErrCommandAlreadyExists = fmt.Errorf("command already exists")
	ErrCommandNotExists     = fmt.Errorf("command not exists")
	ErrContextValueNotFound = fmt.Errorf("context value not found")
	ErrAssetNotFound        = fmt.Errorf("asset not found")
	ErrNotSupported         = fmt.Errorf("not supported")
)
