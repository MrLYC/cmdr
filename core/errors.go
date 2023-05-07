package core

import "fmt"

type ExitError struct {
	code    int
	message string
}

func (e ExitError) Code() int {
	return e.code
}

func (e ExitError) Error() string {
	return fmt.Sprintf("%v: %d", e.message, e.code)
}

func NewExitError(message string, code int) ExitError {
	return ExitError{
		code:    code,
		message: message,
	}
}

var (
	ErrCommandAlreadyActivated = fmt.Errorf("command already activated")
	ErrShellNotSupported       = fmt.Errorf("shell not supported")
	ErrBinaryNotFound          = fmt.Errorf("binaries not found")
	ErrReleaseAssetNotFound    = fmt.Errorf("release asset not found")
)
