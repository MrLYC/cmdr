package utils

import (
	"os"

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

		core.Logger.Error(message, map[string]interface{}{
			"error": err,
		})
		os.Exit(-1)
	}
}

func PanicOnError(message string, errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}

		core.Logger.Error(message, map[string]interface{}{
			"error": err,
		})
		panic(err)
	}
}
