package define

import (
	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"

	"github.com/mrlyc/cmdr/utils"
)

var Logger logur.Logger

func init() {
	Logger = logur.NewNoopLogger()
}

func InitLogger() {
	logrusLogger := logrus.New()
	level, err := logrus.ParseLevel(Configuration.GetString("log.level"))
	utils.CheckError(err)
	logrusLogger.SetLevel(level)

	Logger = logrusadapter.New(logrusLogger)
}
