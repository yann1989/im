// Author       kevin
// Time         2019-08-04 21:25
// File Desc    日志级别

package log

// Level of severity.
type Level int

// Verbose
type Verbose bool

// 常用的一些日志级别
const (
	_debugLevel Level = iota
	_infoLevel
	_warnLevel
	_errorLevel
	_fatalLevel
)

var levelNames = [...]string{
	_debugLevel: "DEBUG",
	_infoLevel:  "INFO",
	_warnLevel:  "WARN",
	_errorLevel: "ERROR",
	_fatalLevel: "FATAL",
}

// String implementation.
func (l Level) String() string {
	return levelNames[l]
}
