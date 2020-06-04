// Author: yann
// Date: 2020/6/3 1:01 下午
// Desc:

package logger

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strconv"
	"sync"
)

//	PanicLevel Level = iota
//	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
//	// logging level is set to Panic.
//	FatalLevel
//	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
//	// Commonly used for hooks to send errors to an error tracking service.
//	ErrorLevel
//	// WarnLevel level. Non-critical entries that deserve eyes.
//	WarnLevel
//	// InfoLevel level. General operational entries about what's going on inside the
//	// application.
//	InfoLevel
//	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
//	DebugLevel
//	// TraceLevel level. Designates finer-grained informational events than the Debug.
//	TraceLevel

type LoggerConfig struct {
	Level int
}

var once = sync.Once{}
var err = "PanicLevel = 0\n" +
	"FatalLevel = 1:\n" +
	"FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the logging level is set to Panic.\n" +
	"ErrorLevel = 2:\n" +
	"ErrorLevel level. Logs. Used for errors that should definitely be noted. Commonly used for hooks to send errors to an error tracking service.\n" +
	"WarnLevel = 3:\n" +
	"WarnLevel level. Non-critical entries that deserve eyes.\n" +
	"InfoLevel = 4:\n" +
	"InfoLevel level. General operational entries about what's going on inside the application.\n" +
	"DebugLevel = 5:\n" +
	"DebugLevel level. Usually only enabled when debugging. Very verbose logging.\n" +
	"TraceLevel = 6:\n" +
	"TraceLevel level. Designates finer-grained informational events than the Debug.\n"

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableTimestamp: false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, filename := path.Split(f.File)
			return f.Function, filename + "/line:" + strconv.Itoa(f.Line)
		},
	})
	//logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
}

func (l *LoggerConfig) verifyParam() {
	if (logrus.Level)(l.Level) < logrus.PanicLevel || (logrus.Level)(l.Level) > logrus.TraceLevel {
		panic(fmt.Sprintf("日志级别配置错误[usage]:\n%s\n", err))
	}
}

func InitLogger(config *LoggerConfig) {
	once.Do(func() {
		config.verifyParam()
		logrus.SetLevel(logrus.Level(config.Level))
		logrus.AddHook(NewFileHook())
	})
}

// 输出log
func Log(request *restful.Request) *logrus.Entry {
	// 输出行数
	//logrus.SetReportCaller(true)
	//// 输出json
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	// 固定输出请求信息
	requestLogger := logrus.WithFields(logrus.Fields{
		"url":    request.Request.RequestURI,
		"method": request.Request.Method,
		"ip":     request.Request.RemoteAddr,
	})
	return requestLogger
}
