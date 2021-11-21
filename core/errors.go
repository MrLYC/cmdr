package core

import (
	"fmt"
)

var (
	ErrCommandAlreadyExists = fmt.Errorf("command already exists")
	ErrCommandNotExists     = fmt.Errorf("command not exists")
	ErrContextValueNotFound = fmt.Errorf("context value not found")
)
