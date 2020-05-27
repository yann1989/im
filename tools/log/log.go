// Author       kevin
// Time         2019-08-04 20:46
// File Desc    日志主文件

package log

import (
	"context"
	"fmt"
)

var (
	h              Handler       // 日志handler
	c              *Config       // 日志配置类
	_defaultStdout = NewStdout() // 日志包默认的日志handler
)

// TODO 可以进一步将rsyslog相关的配置信息进行封装
type Config struct {
	ServiceAddress string `toml:"-"`              // 发送日志的服务的ip+port TODO 可以进一步封装
	RsyslogAddress string `toml:"rsyslogAddress"` // 日志转发服务器的ip+端口号
	Stdout         bool   `toml:"stdout"`         // 是否输出到stdout
}

// D represents a map of entry level data used for structured logging.
// type D map[string]interface{}
type D struct {
	Key   string
	Value interface{}
}

// log package default init
func init() {
	// 默认日志配置
	c = &Config{
		Stdout: true,
	}
	// 默认日志handler
	h = NewStdout()
}

// Init log包自定义初始化, 调用者通过传入 conf 配置信息, 可以自定义日志
// 目前支持两种logger: 1. 打印到控制台; 2.发送给一个日志转发服务器
func Init(conf *Config) {
	*c = *conf
	if !conf.Stdout { // 日志转发
		h = NewRsyslog()
	}
}

// Close 释放log占用的所有资源, 比如打开的连接或文件等,
// 具体释放什么样的资源又 handler的实现类的Close方法决定
func Close() (err error) {
	// 调用具体实现类
	err = h.Close()
	// 场景: 需要自定义一个日志, 但是只需要使用一次, 那么后续可以 defer log.Close()
	// 来关闭自定义的日志. 然后, 恢复回默认的handler, 进程内其他地方调用日志时使用的
	// 就是默认日志handler.
	h = _defaultStdout
	return
}

// KV return a log kv for logging field.
func KV(key string, value interface{}) D {
	return D{
		Key:   key,
		Value: value,
	}
}

// Info logs a message at the info log level.
func Info(format string, args ...interface{}) {
	h.Log(context.Background(), _infoLevel, KV(_msg, fmt.Sprintf(format, args...)))
}

// Warn logs a message at the warning log level.
func Warn(format string, args ...interface{}) {
	h.Log(context.Background(), _warnLevel, KV(_msg, fmt.Sprintf(format, args...)))
}

// Error logs a message at the error log level.
func Error(format string, args ...interface{}) {
	h.Log(context.Background(), _errorLevel, KV(_msg, fmt.Sprintf(format, args...)))
}

// Infov logs a message at the info log level.
func Infov(ctx context.Context, args ...D) {
	h.Log(ctx, _infoLevel, args...)
}

// Warnv logs a message at the warning log level.
func Warnv(ctx context.Context, args ...D) {
	h.Log(ctx, _warnLevel, args...)
}

// Errorv logs a message at the error log level.
func Errorv(ctx context.Context, args ...D) {
	h.Log(ctx, _errorLevel, args...)
}
