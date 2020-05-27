// Author: yann
// Date: 2020/5/23 10:10 上午
// Desc:

package manager

const (
	TOKEN_FIELD      = "authorization"
	NODE_DATA_LENGTH = 64    //每个连接的消息队列长度
	UID              = "uid" //获取用户id常量
	HEART_CHAN       = 1     //心跳chan的值, 用于通知心跳协程
)
