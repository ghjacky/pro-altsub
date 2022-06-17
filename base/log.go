package base

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _log = logrus.New()

func initLog() {
	level, err := logrus.ParseLevel(Config.MainConfig.Level)
	if err != nil {
		_log.Warnln("failed to parse log level, using default value: [INFO]")
		_log.SetLevel(logrus.InfoLevel)
	} else {
		_log.SetLevel(level)
	}
	_log.SetOutput(&lumberjack.Logger{
		Filename: Config.MainConfig.Log,
		MaxSize:  Config.MainConfig.LogRotateMegaBytes, // megabytes
		MaxAge:   Config.MainConfig.LogRetainDays,      //days
		Compress: true,                                 // disabled by default
	})
}

func NewLog(level string, err error, message string, caller string) {
	var f = func() (string) {
		if err != nil {
			return fmt.Sprintf("%s - %s: %s", caller, message, err.Error())
		} else {
			return fmt.Sprintf("%s - %s", caller, message)
		}
	}
	switch level {
	case "debug":
		_log.Debugln(f())
	case "trace":
		_log.Traceln(f())
	case "info":
		_log.Infoln(f())
	case "warn":
		_log.Warnln(f())
	case "error":
		_log.Errorln(f())
	case "fatal":
		_log.Fatalln(f())
	case "panic":
		_log.Panicln(f())
	default:
		if err != nil {
			_log.Errorln(f())
		} else {
			_log.Infoln(f())
		}
	}
}
