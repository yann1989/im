// Author: yann
// Date: 2020/5/23 12:53 下午
// Desc:

package common

type MsgType int

const (
	MSG_TYPE_SINGLE        = iota + 1 //1
	MSG_TYPE_GROUP                    //2
	MSG_TYPE_SINGLE_SYSTEM            //3
	MSG_TYPE_GLOBAL_SYSTEM            //4  全局系统通知
	MSG_TYPE_SERVICE                  //5  客服消息
	MSG_TYPE_HEART         = 99
)

const (
	USER_ONLINE = iota + 1
)
