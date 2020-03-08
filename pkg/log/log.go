package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/sirupsen/logrus"
)

const (
	//包名变化记得修改这里的值
	logrusPkgName      = "github.com/sirupsen/logrus"
	myLoggerPkgName    = "github.com/jinlingan/gringotts/pkg/log"
	functionNamePrefix = "github.com/jinlingan/gringotts/"

	buildDir              = string(os.PathSeparator) + "gringotts" + string(os.PathSeparator)
	buildDirLen           = len(buildDir)
	functionNamePrefixLen = len(functionNamePrefix)
)

// NewStdAndFileLogger 新建 AgentLogger
func NewStdAndFileLogger(filepath string) Logger {
	return &multiLogger{
		parentLogs: []*logrus.Logger{
			{
				Out: &lumberjack.Logger{
					Filename:   filepath,
					MaxSize:    100, // MB
					MaxBackups: 3,
					LocalTime:  true,
					MaxAge:     365, // Days
				},
				Formatter: &logrus.TextFormatter{
					DisableColors:    true,
					CallerPrettyfier: caller,
				},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.DebugLevel,
				ReportCaller: true,
			},
			{
				Out: os.Stdout,
				Formatter: &logrus.TextFormatter{
					FullTimestamp:    true,
					CallerPrettyfier: caller,
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
		parentLogs: []*logrus.Logger{
			{
				Out: os.Stdout,
				Formatter: &logrus.TextFormatter{
					FullTimestamp:    true,
					CallerPrettyfier: caller,
				},
				Hooks:        make(logrus.LevelHooks),
				Level:        logrus.DebugLevel,
				ReportCaller: true,
			},
		},
	}
}

// 目前的实现方式会导致取两次调用栈信息（一次在 logrus 中，一次在这里），
// 如果关闭 logrus 的调用栈会导致 TextFormatter 不调用这个函数
// TODO:看看能不能只调用一次，可能需要自己实现一个 TextFormatter
func caller(_ *runtime.Frame) (function string, file string) {

	pcs := make([]uintptr, 10)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		//pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		pkg := getPackageName(f.Function)
		if pkg != logrusPkgName && pkg != myLoggerPkgName {

			idx := strings.Index(f.File, buildDir)
			filename := f.File
			if idx >= 0 {
				filename = f.File[idx+buildDirLen : len(f.File)]
			}

			file := fmt.Sprintf(" %s:%d", filename, f.Line)
			functionName := f.Function
			fidx := strings.Index(functionName, functionNamePrefix)

			if fidx == 0 {
				functionName = functionName[fidx+functionNamePrefixLen : len(functionName)]
			}

			return functionName, file
		}
	}

	return "unknow", "unknow"
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
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
	parentLogs []*logrus.Logger
}

func (l *multiLogger) Log(level logrus.Level, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Log(level, args...)
	}
}

func (l *multiLogger) Logf(level logrus.Level, format string, args ...interface{}) {
	for _, p := range l.parentLogs {
		p.Logf(level, format, args...)
	}
}

func (l *multiLogger) Loge(level logrus.Level, err error, format string, args ...interface{}) {

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf(format, args...))
	msg.WriteString(fmt.Sprintf(". Caused by: %v\nError Chain is:\n%+v", err, err))

	for _, p := range l.parentLogs {
		p.Log(level, msg.String())
	}
}

//Debugf
func (l *multiLogger) Debugf(format string, args ...interface{}) {
	l.Logf(logrus.DebugLevel, format, args...)
}

//Infof
func (l *multiLogger) Infof(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args...)
}

//Warnf
func (l *multiLogger) Warnf(format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, format, args...)
}

//Errorf
func (l *multiLogger) Errorf(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args...)
}

//Fatalf
func (l *multiLogger) Fatalf(format string, args ...interface{}) {
	l.Logf(logrus.FatalLevel, format, args...)
}

//Debug
func (l *multiLogger) Debug(args ...interface{}) {
	l.Log(logrus.DebugLevel, args...)
}

//Info
func (l *multiLogger) Info(args ...interface{}) {
	l.Log(logrus.InfoLevel, args...)
}

//Warn
func (l *multiLogger) Warn(args ...interface{}) {
	l.Log(logrus.WarnLevel, args...)
}

//Error
func (l *multiLogger) Error(args ...interface{}) {
	l.Log(logrus.ErrorLevel, args...)
}

//Fatal
func (l *multiLogger) Fatal(args ...interface{}) {
	l.Log(logrus.FatalLevel, args...)
}

//Debuge
func (l *multiLogger) Debuge(err error, format string, args ...interface{}) {
	l.Loge(logrus.DebugLevel, err, format, args...)
}

//Infoe
func (l *multiLogger) Infoe(err error, format string, args ...interface{}) {
	l.Loge(logrus.InfoLevel, err, format, args...)
}

//Warne
func (l *multiLogger) Warne(err error, format string, args ...interface{}) {
	l.Loge(logrus.WarnLevel, err, format, args...)
}

//Errore
func (l *multiLogger) Errore(err error, format string, args ...interface{}) {
	l.Loge(logrus.ErrorLevel, err, format, args...)
}

//Fatale
func (l *multiLogger) Fatale(err error, format string, args ...interface{}) {
	l.Loge(logrus.FatalLevel, err, format, args...)
}
