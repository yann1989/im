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
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_SINGLE_PREFIX, message.ContactID),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
		if err = b.sendOrBroadcast(message.ToID); err != nil {
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
		return mq.Broadcast(b.jsonData)
	//存入系统历史消息
	case common.MSG_TYPE_SYSTEM:
		if err = utils.RedisUtils.ZaddSingle(
			utils.RedisUtils.BuildKey(common.REDIS_KEY_HISTORY_MESSAGE_SYSTEM_PREFIX, message.ToID),
			message.CreateTime,
			string(b.jsonData)); err != nil {
			return err
		}
		if err = b.sendOrBroadcast(message.ToID); err != nil {
			return err
		}
	//心跳
	case common.MSG_TYPE_HEART:
	}

	return nil
}

func (b *broadcast) sendOrBroadcast(userID int64) error {
	manager.rwlocker.RLock()
	defer manager.rwlocker.RUnlock()
	if node, has := manager.nodes[userID]; has {
		err := mq.Broadcast(b.jsonData)
		if err != nil {
			return err
		}
		node.Conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
		node.Conn.WriteMessage(websocket.TextMessage, b.jsonData)
		return nil
	}

	return mq.Broadcast(b.jsonData)
}
