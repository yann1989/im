// Author       kevin
// Time         2019-08-04 21:25
// File Desc    日志 stdout handler

package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
	file_log "yann-chat/tools/log/file-log"
)

// StdoutHandler stdout log handler
type StdoutHandler struct {
	out io.Writer
}

// NewStdout create a stdout log handler
func NewStdout() *StdoutHandler {
	return &StdoutHandler{
		out: os.Stdout,
	}
}

// stdout-handler 日志输出方法
// 参数: context, 级别, 日志输出的内容(Key-value格式)
func (h *StdoutHandler) Log(ctx context.Context, lv Level, args ...D) {

	// 存储日志输出内容(Key-value)的map, 比如, msg="hello golang"
	d := make(map[string]interface{}, 10+len(args))
	for _, arg := range args {
		d[arg.Key] = arg.Value
	}

	callerInfo := callerPrettify(getCaller())

	date := time.Now().Format("2006-01-02 15:04:05")
	// 输出格式: "[级别] 位置: 日志内容", 比如, "[INFO] go-taurus/app/user_service/cmd/main.go-main-44: user-dao start"
	output := fmt.Sprintf("[%s][%s] %s: %s\n", lv.String(), date, callerInfo, d[_msg])

	h.out.Write([]byte(output))

	// (输出)写入日志文件---不想写入日志可以不加下面这行代码,和上面的逻辑没有耦合
	file_log.NewRealStLogger(int(lv)).ERROR(output)
}

// Close stdout-handler has nothing to close
func (h *StdoutHandler) Close() error {
	return nil
}
