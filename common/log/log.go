package log

import (
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// NewStdAndFileLogger 新建 AgentLogger
func NewStdAndFileLogger(filepath string) Logger {
	return &multiLogger{
		parentLogs: []logrus.FieldLogger{
			&logrus.Logger{
				Out: &lumberjack.Logger{
					Filename:   filepath,
					MaxSize:    100, // MB
					MaxBackups: 3,
					LocalTime:  true,
					MaxAge:     365, // Days
				},
				Formatter: &logrus.TextFormatter{
					DisableColors: true,
				},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.DebugLevel,
				ReportCaller: true,
			}, &logrus.Logger{
				Out: os.Stdout,
				Formatter: &logrus.TextFormatter{
					FullTimestamp: true,
				},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.DebugLevel,
				ReportCaller: true,
			},
		},
	}
}

// NewStdoutLogger 新建 AgentLogger
func NewStdoutLogger() Logger {
	return &multiLogger{
		parentLogs: []logrus.FieldLogger{
			&logrus.Logger{
				Out: os.Stdout,
				Formatter: &logrus.TextFormatter{
					FullTimestamp: true,
				},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.DebugLevel,
				ReportCaller: true,
			},
		},
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

	Debuge(err error, format string, args ...interface{})
	Infoe(err error, format string, args ...interface{})
	Warne(err error, format string, args ...interface{})
	Errore(err error, format string, args ...interface{})
	Fatale(err error, format string, args ...interface{})
}

type multiLogger struct {
	parentLogs []logrus.FieldLogger
}

//TODO 剥离循环
//TODO 打印 error 具体堆栈
//Debugf
func (l *multiLogger) Debugf(format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Debugf(format, args...)
	}
}

//Infof
func (l *multiLogger) Infof(format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Infof(format, args...)
	}
}

//Warnf
func (l *multiLogger) Warnf(format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Warnf(format, args...)
	}
}

//Errorf
func (l *multiLogger) Errorf(format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Errorf(format, args...)
	}
}

//Fatalf
func (l *multiLogger) Fatalf(format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Fatalf(format, args...)
	}
}

//Debug
func (l *multiLogger) Debug(args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Debug(args...)
	}
}

//Info
func (l *multiLogger) Info(args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Info(args...)
	}
}

//Warn
func (l *multiLogger) Warn(args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Warn(args...)
	}
}

//Error
func (l *multiLogger) Error(args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Error(args...)
	}
}

//Fatal
func (l *multiLogger) Fatal(args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Fatal(args...)
	}
}

//Debuge
func (l *multiLogger) Debuge(err error, format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Debug(errors.Wrapf(err, format, args...))
	}
}

//Infoe
func (l *multiLogger) Infoe(err error, format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Info(errors.Wrapf(err, format, args...))
	}
}

//Warne
func (l *multiLogger) Warne(err error, format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Warn(errors.Wrapf(err, format, args...))
	}
}

//Errore
func (l *multiLogger) Errore(err error, format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Error(errors.Wrapf(err, format, args...))
	}
}

//Fatale
func (l *multiLogger) Fatale(err error, format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Fatal(errors.Wrapf(err, format, args...))
	}
}
