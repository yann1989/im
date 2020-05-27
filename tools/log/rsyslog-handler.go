// Author       kevin
// Time         2019-08-04 21:42
// File Desc    Rsyslog 日志handler

package log

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
)

// StdoutHandler stdout log handler
type RsyslogHandler struct {
	out io.Writer
}

// NewStdout create a stdout log handler
func NewRsyslog() *RsyslogHandler {
	// 日志输出到日志转发服务器
	conn, err := net.Dial("tcp", c.RsyslogAddress)
	if err != nil {
		log.Fatal("Unable to connect rsyslog")
		os.Exit(1)
	}
	return &RsyslogHandler{
		out: conn,
	}
}

// Log 日志转发, 格式JSON
func (h *RsyslogHandler) Log(ctx context.Context, lv Level, args ...D) {
	d := make(map[string]interface{}, 10+len(args))
	for _, arg := range args {
		d[arg.Key] = arg.Value
	}
	// 添加公共的日志信息
	addExtraFields(d)
	d[_level] = lv.String()
	jsonString, _ := json.Marshal(d)
	h.out.Write(jsonString)
	h.out.Write([]byte("\n"))
}

// Close 关闭handler
// FIXME: 需要关闭在创建 rsyslogHandler 时和日志转发服务器之间建立的连接
func (h *RsyslogHandler) Close() error {
	return nil
}
