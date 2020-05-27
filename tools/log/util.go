// Author       kevin
// Time         2019-08-04 20:46
// File Desc    helper functions for log package

package log

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	maximumCallerDepth int = 25
	knownLogFrames     int = 4 // 从进入log包, 到调用当前文件, 经过的frame的个数(底层对应runtime.Caller()的参数skip number)
)

var (
	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// qualified package name, cached at first use
	logPackage string
	// Used for caller information initialisation
	callerInitOnce sync.Once
)

// 添加日志输出的公共信息
func addExtraFields(fields map[string]interface{}) {
	// 时间戳
	fields[_time] = time.Now().Format(_timeFormat)
	// 服务ip
	fields[_serviceIP] = c.ServiceAddress
	// 日志调用者信息
	callerInfo := callerPrettify(getCaller())
	fields[_caller] = callerInfo
}

// callerPrettify 调用者信息格式美化, packageName/functionName.go:line
func callerPrettify(caller *runtime.Frame) (callerInfo string) {

	// file
	fileName := caller.File
	index := strings.LastIndex(fileName, "yann-chat")
	fileName = fileName[index:]

	// function
	s := strings.Split(caller.Function, ".")
	funcName := s[len(s)-1]

	// line
	lineNumber := caller.Line

	callerInfo = fileName + "-" + funcName + "-" + strconv.Itoa(lineNumber)
	return
}

// getCaller retrieves the name of the first non-log package calling function
func getCaller() *runtime.Frame {

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		logPackage = getPackageName(runtime.FuncForPC(pcs[1]).Name())

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary
		minimumCallerDepth = knownLogFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logPackage {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
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
