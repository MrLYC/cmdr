package core

import "errors"

var (
	ErrCommandAlreadyExists = errors.New("command already exists")
	ErrCommandNotExists     = errors.New("command not exists")
)
