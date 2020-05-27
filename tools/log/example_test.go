// Author       kevin
// Time         2019-08-06 13:12
// File Desc    日志示例

package log_test

import "taurus-admin/library/log"

// 默认的log, 输出到控制台, 直接import log即可
func ExampleDefaultLog() {
	msg := "hello golang"
	log.Error("%s", msg)
	log.Warn("%s", msg)
	log.Info("%s", msg)
	// output:
}

// 使用自定义log
func ExampleCustomizeLog() {

	// 日志配置信息, 一般通过加载文件来初始化
	config := new(log.Config)

	// 调用日志的进程所在的机器的ip地址+端口号
	config.ServiceAddress = "192.111.234.1:9089"

	// 日志是否打印到stdout, 如果为true, 下面关于rsyslog转发的配置将被忽视
	config.Stdout = true

	// 日志转发服务器ip
	config.RsyslogAddress = "127.0.0.1:514"

	// 初始化自定义的日志
	log.Init(config)

	// 使用
	msg := "hello golang"
	log.Error("%s", msg)
	log.Warn("%s", msg)
	log.Info("%s", msg)

	// output:
}
