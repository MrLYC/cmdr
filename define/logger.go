package define

import (
	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
)

var Logger logur.Logger

func init() {
	Logger = logur.NewNoopLogger()
}

func InitLogger() {
	logrusLogger := logrus.New()
	level, err := logrus.ParseLevel(Configuration.GetString("log.level"))
	if err != nil {
		panic(err)
	}
	logrusLogger.SetLevel(level)

	Logger = logrusadapter.New(logrusLogger)
}
