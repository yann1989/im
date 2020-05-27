// Author: yann
// Date: 2020/5/23 3:01 下午
// Desc:

package manager

import (
	"time"
	"yann-chat/tools/log"
)

func Receiver(node *Node) {
	log.Info("读事件开始读取数据")
	err := node.Conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if err != nil {
		return
	}
	_, data, err := node.Conn.ReadMessage()
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("读取到消息: %s", string(data))

	//node.HeartCh <- HEART_CHAN
	//把消息广播到局域网
	if err = NewBroadcaster(data).Execute(); err != nil {
		log.Error("广播消息失败: %s", err.Error())
	}
	return
}
