package base

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type logger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
	lg                    *logrus.Logger
}

func newGormLog(l *logrus.Logger) *logger {
	return &logger{
		SkipErrRecordNotFound: false,
		Debug:                 true,
		lg:                    l,
	}
}

func (l *logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.lg.WithContext(ctx).Infof(s, args)
}

func (l *logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.lg.WithContext(ctx).Warnf(s, args)
}

func (l *logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.lg.WithContext(ctx).Errorf(s, args)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.lg.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.lg.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	if l.Debug {
		l.lg.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
	}
}

var _log = logrus.New()
var gormLog *logger

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
	gormLog = newGormLog(_log)
	if level >= logrus.DebugLevel {
		gormLog.Debug = true
	} else {
		gormLog.Debug = false
	}
}

func NewLog(level string, err error, message string, caller string) {
	var f = func() string {
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
