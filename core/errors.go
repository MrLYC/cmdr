package core

import "fmt"

type ExitError struct {
	code int
}

func (e ExitError) Code() int {
	return e.code
}

func (e ExitError) Error() string {
	return fmt.Sprintf("exit code %d", e.code)
}

func NewExitError(code int) ExitError {
	return ExitError{
		code: code,
	}
}

var (
	ErrCommandAlreadyActivated = fmt.Errorf("command already activated")
	ErrShellNotSupported       = fmt.Errorf("shell not supported")
	ErrBinaryNotFound          = fmt.Errorf("binaries not found")
	ErrReleaseAssetNotFound    = fmt.Errorf("release asset not found")
)
