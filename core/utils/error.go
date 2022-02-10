package utils

import (
	"fmt"
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

func ExitWithError(err error, message string, args ...interface{}) {
	if err == nil {
		return
	}

	core.Logger.Error(fmt.Sprintf(message, args...), map[string]interface{}{
		"error": err,
	})
	os.Exit(-1)
}
