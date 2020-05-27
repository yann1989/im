// Author: yann
// Date: 2020/5/23 9:33 上午
// Desc:

package manager

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"
	"net/http"
	"yann-chat/common"
	"yann-chat/tools/log"
	"yann-chat/tools/utils"
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
		log.Error("用户 %d 建立websocket失败: ", uid)
		return nil
	}

	//将连接放入管理, 添加失败则关闭连接返回
	node := manager.Add(uid, conn)
	if node == nil {
		conn.Close()
		return nil
	}

	//设置在线状态 失败则删除管理中的节点
	if err = utils.RedisUtils.Hset(string(common.REDIS_KEY_USER_ONLINE_PREFIX), uid, common.USER_ONLINE); err != nil {
		log.Error("redis写入失败")
		manager.Remove(uid)
		return nil
	}

	return node
}
