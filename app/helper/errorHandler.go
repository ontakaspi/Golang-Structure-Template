package helper

import (
	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
	"golang-example/libraries/logger"
)

func CatchError() {
	if err := recover(); err != nil {

		logger.SetLogFileAndConsole(logger.LogData{
			Message: "unexpected error",
			CustomFields: logrus.Fields{
				"message": err,
			},
			Level: "ERROR",
		})

		dataErrr := tracerr.Wrap(err.(error))
		tracerr.PrintSourceColor(dataErrr)
	}
}

func ErrorHandler(err error) {
	panic(err)
}
