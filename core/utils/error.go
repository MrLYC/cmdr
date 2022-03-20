package utils

import (
	"github.com/mrlyc/cmdr/core"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func CallClose(closer interface {
	Close() error
}) {
	CheckError(closer.Close())
}

func ExitOnError(message string, errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}

		core.GetLogger().Error(message, map[string]interface{}{
			"error": err,
		})
		panic(core.NewExitError(-1))
	}
}

func PanicOnError(message string, errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}

		core.GetLogger().Error(message, map[string]interface{}{
			"error": err,
		})
		panic(err)
	}
}
