// Author: yann
// Date: 2020/5/23 11:18 上午
// Desc:

package manager

import (
	"yann-chat/common"
	"yann-chat/model"
	"yann-chat/tools/log"
	"yann-chat/tools/mq"
	"yann-chat/tools/snowflake"
	"yann-chat/tools/utils"
)

type broadcast struct {
	jsonData []byte
}

//构造广播
func NewBroadcaster(params []byte) *broadcast {
	broadcaster := new(broadcast)
	broadcaster.jsonData = params
	return broadcaster
}

//执行
func (b *broadcast) Execute() error {
	//心跳
	if len(b.jsonData) == 0 {
		return nil
	}
	//获取消息
	var err error
	message := new(model.Message)
	if err = utils.JsonUnMarshal(b.jsonData, message); err != nil {
		return common.ERR_BAD_REQUES
	}
	if message.MsgType == 0 || message.CreateTime == 0 {
		return common.ERR_BAD_REQUES
	}
	//存入redis历史消息
	message.MsgID = snowflake.NextId()
	switch message.MsgType {
	//存入私聊历史消息
	case common.MSG_TYPE_SINGLE:
		if err = utils.RedisUtils.ZaddSingle(
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_SINGLE_PREFIX, message.ContactID),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
	//存入群聊历史消息
	case common.MSG_TYPE_GROUP:
		if err = utils.RedisUtils.ZaddSingle(
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_GROUP_PREFIX, message.ToID),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
	//存入系统历史消息
	case common.MSG_TYPE_SYSTEM:
		if err = utils.RedisUtils.ZaddSingle(
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_SYSTEM_PREFIX, message.ToID),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
		////存入客服历史消息  todo 新增连接类型 type 1=用户 2=客服
		//case common.MSG_TYPE_SERVICE:
		//	if err = utils.RedisUtils.ZaddSingle(
		//		utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_SERVICE_PREFIX, uid),
		//		message.CreateTime,
		//		string(b.jsonData)); err != nil {
		//		return err
		//	}
	case common.MSG_TYPE_HEART:
		return nil
	}

	//广播到所有节点 todo 判断如果为本机节点则直接发送, 否则广播到所有节点
	log.Info("历史消息写入redis成功")
	err = mq.Broadcast(b.jsonData)
	if err != nil {
		return err
	}
	log.Info("消息广播成功")
	return nil
}
