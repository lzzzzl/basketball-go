package log

import "github.com/sirupsen/logrus"

var log *logrus.Logger

func LogInit(logLevel string) {
	log = logrus.New()

	level := logrus.DebugLevel
	switch {
	case logLevel == "debug":
		level = logrus.DebugLevel
	case logLevel == "info":
		level = logrus.InfoLevel
	case logLevel == "error":
		level = logrus.ErrorLevel
	case logLevel == "warn":
		level = logrus.WarnLevel
	default:
		level = logrus.DebugLevel
	}
	log.Formatter = &logrus.JSONFormatter{}
	log.SetLevel(level)
}

func Println(v ...interface{}) {
	log.Info(v)
}

func Error(v ...interface{}) {
	log.Error(v)
}

func Debug(v ...interface{}) {
	log.Debug(v)
}
