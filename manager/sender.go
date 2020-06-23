// Author: yann
// Date: 2020/5/23 2:58 下午
// Desc:

package manager

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
	"yann-chat/common"
	"yann-chat/model"
	"yann-chat/tools/dao/redisClient"

	"yann-chat/tools/utils"
)

//消费redis subscription
func (m *ConnectManager) startConsume() {
	for ch := range redisClient.Ch {
		manager.gopool.Schedule(func() {
			m.dispatch([]byte(ch.Payload))
		})
	}
	logrus.Infof("退出conusme")
}

//分发消息
func (m *ConnectManager) dispatch(data []byte) {
	// 解析data为message
	msg := new(model.Message)
	err := json.Unmarshal(data, msg)
	if err != nil {
		logrus.Errorf("消息格式有误, json序列化失败: %s", err.Error())
		return
	}

	//根据消息类型分发消息
	switch msg.MsgType {
	//私聊, 定点系统通知, 客服消息
	case common.MSG_TYPE_SINGLE, common.MSG_TYPE_SINGLE_SYSTEM, common.MSG_TYPE_SERVICE:
		//查看目标用户是否在线,在线则发送
		manager.rwlocker.RLock()
		defer manager.rwlocker.RUnlock()
		node, ok := manager.nodes[msg.ToID]
		if ok {
			node.Conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
			node.Conn.WriteMessage(websocket.TextMessage, data)
		}

	//群聊  可以通过redis获取群成员(遍历群成员) 也可以通过在用户建立连接的时候每个节点上维护一个此用户加入的所有群聊(遍历所有节点)
	case common.MSG_TYPE_GROUP:
		//1.获取群成员
		list, err := utils.RedisUtils.SmembersInt64(utils.RedisUtils.BuildKey(common.REDIS_KEY_GROUP_MEMBER_PREFIX, msg.ToID))
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
		if list == nil || len(list) <= 0 {
			return
		}

		//2.如果群成员在线则发送数据
		manager.rwlocker.RLock()
		defer manager.rwlocker.RUnlock()
		for _, uid := range list {
			node, ok := manager.nodes[uid]
			if ok {
				node.Conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
				node.Conn.WriteMessage(websocket.TextMessage, data)
			}
		}
	//全局系统通知
	case common.MSG_TYPE_GLOBAL_SYSTEM:
		//发送给所有人.
		manager.rwlocker.RLock()
		defer manager.rwlocker.RUnlock()
		for _, node := range manager.nodes {
			node.Conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
			node.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
