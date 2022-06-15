package base

import (
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
	var f = func() (string, string, string, string) {
		return "%s - %s: %s", caller, message, err.Error()
	}
	switch level {
	case "debug":
		_log.Debugf(f())
	case "trace":
		_log.Tracef(f())
	case "info":
		_log.Infof(f())
	case "warn":
		_log.Warnf(f())
	case "error":
		_log.Errorf(f())
	case "fatal":
		_log.Fatalf(f())
	case "panic":
		_log.Panicf(f())
	default:
		if err != nil {
			_log.Errorf(f())
		} else {
			_log.Infof(f())
		}
	}
}
