// Author       kevin
// Time         2019-08-08 20:16
// File Desc    路由表

package main

import (
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"yann-chat/tools/view"
)

// 初始化ws的路由表
func (h *HttpServer) initRoutes(ws *restful.WebService) {
	// websocket
	ws.Route(ws.GET("/websocket").Filter(h.testTokenFilter).To(panicHandle(h.websocketHandel)))
}

//拦截请求中的未知panic
func panicHandle(handle restful.RouteFunction) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("panic recover: %v请求=>%v", req.Request.URL.Path, err)
				view.Response500().ReturnResult(req, res)
			}
		}()
		handle(req, res)
	}
}
