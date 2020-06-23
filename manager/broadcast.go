// Author: yann
// Date: 2020/5/23 11:18 上午
// Desc:

package manager

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
	"yann-chat/common"
	"yann-chat/model"
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
	if err = json.Unmarshal(b.jsonData, message); err != nil {
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
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_SINGLE_PREFIX, message.SessionID),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
		return b.sendOrBroadcast(message.ToID)
	//存入群聊历史消息
	case common.MSG_TYPE_GROUP:
		if err = utils.RedisUtils.ZaddSingle(
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_GROUP_PREFIX, message.ToID),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
		//return mq.Broadcast(b.jsonData)
		return utils.RedisUtils.Broadcast(b.jsonData)
	//存入定点系统历史消息
	case common.MSG_TYPE_SINGLE_SYSTEM:
		if err = utils.RedisUtils.ZaddSingle(
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_SINGLE_SYSTEM_PREFIX, message.ToID),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
		return b.sendOrBroadcast(message.ToID)
	//全局系统通知
	case common.MSG_TYPE_GLOBAL_SYSTEM:
		if err = utils.RedisUtils.ZaddSingle(
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_GLOBAL_SYSTEM, nil),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
		//return mq.Broadcast(b.jsonData)
		return utils.RedisUtils.Broadcast(b.jsonData)
	//todo 客服消息
	case common.MSG_TYPE_SERVICE:
		//if err = utils.RedisUtils.ZaddSingle(
		//	utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_GLOBAL_SYSTEM, nil),
		//	message.CreateTime,
		//	string(b.jsonData)); err != nil {
		//	return err
		//}
		return b.sendOrBroadcast(message.ToID)
	//心跳
	case common.MSG_TYPE_HEART:
	}

	return nil
}

//***************************************************
//Description : 私聊类型才调用, 如果在当前节点则发送, 否则广播
//param : 用户id
//return : 错误信息
//***************************************************
func (b *broadcast) sendOrBroadcast(userID int64) error {
	manager.rwlocker.RLock()
	defer manager.rwlocker.RUnlock()
	if node, has := manager.nodes[userID]; has {
		node.Conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
		node.Conn.WriteMessage(websocket.TextMessage, b.jsonData)
		//manager.dispatch(b.jsonData)
		return nil
	}
	//return mq.Broadcast(b.jsonData)
	return utils.RedisUtils.Broadcast(b.jsonData)
}
