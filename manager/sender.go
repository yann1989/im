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
	"yann-chat/tools/mq"
	"yann-chat/tools/utils"
)

//消费mq
func (m *ConnectManager) startConsume() {
	msgs, err := mq.StartConsume()
	if err != nil {
		logrus.Errorf("amqp 开始消费失败, 失败原因:%s", err.Error())
		panic("amqp 开始消费失败, 失败原因:%s" + err.Error())
	}
	for msg := range msgs {
		manager.gopool.Schedule(func() {
			m.dispatch(msg.Body)
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

	//私聊
	case common.MSG_TYPE_SINGLE, common.MSG_TYPE_SYSTEM, common.MSG_TYPE_SERVICE:
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
		list := utils.RedisUtils.SmembersInt64(utils.RedisUtils.BuildKey(common.REDIS_KEY_GROUP_MEMBER_PREFIX, msg.ToID))
		if list == nil {
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
	}
}
