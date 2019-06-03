package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewStdAndFileLogger 新建 AgentLogger
func NewStdAndFileLogger(filepath string) Logger {
	return &logrus.Logger{
		Out: io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   filepath,
			MaxSize:    100, // MB
			MaxBackups: 3,
			LocalTime:  true,
			MaxAge:     365, // Days
		}),
		Formatter: &logrus.TextFormatter{
			DisableColors: true,
		},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ReportCaller: true,
	}
}

// Logger 日志接口
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// NewStdoutLogger 新建 AgentLogger
func NewStdoutLogger() Logger {
	return &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			DisableColors: true,
		},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ReportCaller: true,
	}
}
