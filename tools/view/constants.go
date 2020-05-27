// Author       yann
// Time         2020-03-10 16:28
// File Desc    请求头, 状态码常量

package view

// 公共错误码
const (
	CODE_SERVER_ERR = 100500 //服务异常
	CODE_TOKEN_ERR  = 100401
)

const (
	ERR_SERVER = "Server is busy, please try again later."
	ERR_TOKEN  = "Identity verification failed, please log in again."
)
