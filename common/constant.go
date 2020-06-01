// Author: yann
// Date: 2020/5/23 12:53 下午
// Desc:

package common

type MsgType int

const (
	MSG_TYPE_SINGLE  = iota + 1 //1
	MSG_TYPE_GROUP              //2
	MSG_TYPE_SYSTEM             //3
	MSG_TYPE_SERVICE            //4
	MSG_TYPE_HEART   = 99
)

const (
	USER_ONLINE = iota + 1
)
