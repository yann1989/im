// Author: yann
// Date: 2020/6/3 1:09 下午
// Desc:

package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	file_log "yann-chat/tools/logger/file-log"
)

type FileHook struct{}

func NewFileHook() *FileHook {
	return new(FileHook)
}

func (f FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f FileHook) Fire(entry *logrus.Entry) error {
	// 写入文件
	msg, _ := entry.String()
	file_log.NewRealStLogger(cast.ToInt(logrus.ErrorLevel)).ERROR(msg)

	return nil
}
