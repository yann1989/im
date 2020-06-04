// Author: yann
// Date: 2020/5/23 9:33 上午
// Desc:

package manager

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"
	"github.com/sirupsen/logrus"
	"net/http"
	"yann-chat/tools/view"
)

type Node struct {
	ClientId int64
	Conn     *websocket.Conn
	fd       *netpoll.Desc
}

//***************************************************
//Description : 同意连接请求
//param :       request
//param :       response
//return :      成功返回连接节点 失败返回nil
//***************************************************
func Accept(request *restful.Request, response *restful.Response) *Node {
	//获取用户id
	uid, ok := request.Attribute(UID).(int64)
	if !ok || uid <= 0 {
		return nil
	}

	//是否超过最大连接限制
	if manager.isMax() {
		return nil
	}

	//创建websocket连接
	conn, err := (&websocket.Upgrader{
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			bytes, _ := json.Marshal(view.Response500())
			w.Write(bytes)
		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(response.ResponseWriter, request.Request, nil)
	if err != nil {
		logrus.Errorf("用户 %d 建立websocket失败: ", uid)
		return nil
	}

	//将连接放入管理, 添加失败则关闭连接返回
	node := manager.Add(uid, conn)
	if node == nil {
		conn.Close()
		return nil
	}

	//todo 根据业务需求, 可记录在线状态

	return node
}
