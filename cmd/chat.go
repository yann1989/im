// Author: yann
// Date: 2019/12/6 下午1:39
// Desc: websocket 通信

package main

import (
	"github.com/emicklei/go-restful"
	"yann-chat/manager"
)

func (HttpServer) websocketHandel(request *restful.Request, response *restful.Response) {

	node := manager.Accept(request, response)
	if node == nil {
		return
	}

	//go manager.HeartBeating(node) //心跳协程
}
