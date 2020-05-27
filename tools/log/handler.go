// Author       kevin
// Time         2019-08-04 21:22
// File Desc    日志handler, 用来处理日志事件, 不同的handler对应不同的功能,
// 				比如, RsyslogHandler主要是用来处理日志转发rsyslog的.

package log

import "context"

const (
	_timeFormat = "2006-01-02T15:04:05.999999"

	// log level defined in level.go.
	_levelValue = "level_value"
	//  log level name: INFO, WARN...
	_level = "level"
	// log time.
	_time = "time"
	// request path.
	// _title = "title"
	// log file.
	_source = "source"
	// 日志信息
	_msg = "msg"
	// app name.
	_appID = "app_id"
	// Service ip address
	_serviceIP = "ip"
	// container ID.
	_instanceID = "instance_id"
	// uniq ID from trace.
	_tid = "traceid"
	// request time.
	// _ts = "ts"
	// 日志调用者
	_caller = "caller"
	// container environment: prod, pre, uat, fat.
	_deplyEnv = "env"
	// cluster.
	_cluster = "cluster"
)

// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services. See the "handlers"
// directory for implementations.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	// Log handle log
	// variadic D is k-v struct represent log content
	Log(context.Context, Level, ...D)

	// 设置日志的输出格式, todo
	// SetFormat set render format on log output
	// see StdoutHandler.SetFormat for detail
	// SetFormat(string)

	// Close handler
	Close() error
}
