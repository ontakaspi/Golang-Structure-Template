package logger

import (
	"github.com/sirupsen/logrus"
	"golang-example/config"
)

// LogData is a struct that contains the data for the logger
type LogData struct {
	Message      string
	Level        string
	CustomFields logrus.Fields
}

// SetLogConsole is a wrapper for logrus.Logger. It is used to log messages to stdout.
// Level sets the log level, by default it is set to info.
// LogMessage logs a message to the log stdout.
// CustomFields is a map of custom fields that can be added to the log message. (OPTIONAL)
// instantiate your own:
//    Logger.SetLog{
//      Message: "string message",
//      Level:  DEBUG|INFO|WARN|ERROR|FATAL,
//      CustomFields: map[string]interface{}{"key": "value"},
//    }
//
func SetLogConsole(LogData LogData) {

	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}

	logger.SetFormatter(formatter)

	switch LogData.Level {
	case "DEBUG":
		logger.WithFields(LogData.CustomFields).Debug(LogData.Message)
	case "WARN":
		logger.WithFields(LogData.CustomFields).Warn(LogData.Message)
	case "ERROR":
		logger.WithFields(LogData.CustomFields).Error(LogData.Message)
	case "FATAL":
		logger.WithFields(LogData.CustomFields).Error(LogData.Message)
	case "PANIC":
		logger.WithFields(LogData.CustomFields).Error(LogData.Message)
	case "TRACE":
		logger.WithFields(LogData.CustomFields).Trace(LogData.Message)
	default:
		logger.WithFields(LogData.CustomFields).Info(LogData.Message)
	}

	return
}

// SetLogFile is a wrapper for logrus.Logger. It is used to log messages to file.
// Level sets the log level, by default it is set to info.
// LogMessage logs a message to the log file.
// CustomFields is a map of custom fields that can be added to the log message. (OPTIONAL)
// instantiate your own:
//    Logger.SetLog{
//      Message: "string message",
//      Level:  DEBUG|INFO|WARN|ERROR|FATAL,
//      CustomFields: map[string]interface{}{"key": "value"},
//    }
//
func SetLogFile(LogData LogData) {

	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}

	logger.SetFormatter(formatter)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(config.LogFile)

	switch LogData.Level {
	case "DEBUG":
		logger.WithFields(LogData.CustomFields).Debug(LogData.Message)
	case "WARN":
		logger.WithFields(LogData.CustomFields).Warn(LogData.Message)
	case "ERROR":
		logger.WithFields(LogData.CustomFields).Error(LogData.Message)
	case "FATAL":
		logger.WithFields(LogData.CustomFields).Fatal(LogData.Message)
	case "PANIC":
		logger.WithFields(LogData.CustomFields).Panic(LogData.Message)
	case "TRACE":
		logger.WithFields(LogData.CustomFields).Trace(LogData.Message)
	default:
		logger.WithFields(LogData.CustomFields).Info(LogData.Message)
	}

	return
}

// SetLogFileAndConsole is a wrapper for logrus.Logger. It is used to log messages to file and stdout.
func SetLogFileAndConsole(LogData LogData) {

	SetLogConsole(LogData)
	SetLogFile(LogData)
	return
}
