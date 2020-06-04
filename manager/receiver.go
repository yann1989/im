// Author: yann
// Date: 2020/5/23 3:01 下午
// Desc:

package manager

import (
	"github.com/sirupsen/logrus"
	"time"
)

func Receiver(node *Node) {
	err := node.Conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if err != nil {
		return
	}
	_, data, err := node.Conn.ReadMessage()
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
	logrus.Infof("读取到消息: %s", string(data))

	//node.HeartCh <- HEART_CHAN
	//把消息广播到局域网
	if err = NewBroadcaster(data).Execute(); err != nil {
		logrus.Errorf("广播消息失败: %s", err.Error())
	}
	return
}
