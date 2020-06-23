// Author: yann
// Date: 2020/5/23 12:48 下午
// Desc:

package common

type RedisKey string

const (
	REDIS_KEY_HISTORY_MESSAGE_SINGLE_PREFIX        RedisKey = "history:message:single:"        //后面跟会话id history:message:single:1231231223   私聊历史消息
	REDIS_KEY_HISTORY_MESSAGE_GROUP_PREFIX         RedisKey = "history:message:group:"         //后面跟群聊id history:message:group:1231231223    群聊历史消息
	REDIS_KEY_HISTORY_MESSAGE_SINGLE_SYSTEM_PREFIX RedisKey = "history:message:single:system:" //后面跟用户id history:message:single:system:1231231223   系统历史消息
	REDIS_KEY_HISTORY_MESSAGE_GLOBAL_SYSTEM        RedisKey = "history:message:global:system"  //后面跟用户id history:message:global:system  // 全局系统通知
	REDIS_KEY_GROUP_MEMBER_PREFIX                  RedisKey = "group:member:"                  //后面跟群id group:member:1231231223   //群成员
)
