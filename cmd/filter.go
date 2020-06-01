// Author: yann
// Date: 2019/12/14 下午1:39
// Desc: 请求拦截器

package main

import (
	"github.com/emicklei/go-restful"
	"github.com/spf13/cast"
	"time"
	"yann-chat/manager"
	"yann-chat/tools/jwt"
	"yann-chat/tools/log"
	"yann-chat/tools/view"
)

func (h *HttpServer) tokenFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	//校验token
	jwtToken, err := jwt.VerifyAndRenewToken(req.QueryParameter(manager.TOKEN_FIELD), h.privateKey)
	if err != nil {
		view.Response401(req).ReturnResult(req, resp)
		return
	}
	// 校验成功，继续
	log.Info("时间:%s--请求来自:%s--用户ID:%d", time.Now().String(), req.Request.RemoteAddr, jwtToken.Claims.UserId)
	req.SetAttribute("uid", jwtToken.Claims.UserId)
	chain.ProcessFilter(req, resp)
}

//测试用, 直接传websocket?uid=xxxxx
func (h *HttpServer) testTokenFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	//校验token
	uid := cast.ToInt64(req.QueryParameter(manager.UID))
	if uid == 0 {
		view.Response401(req).ReturnResult(req, resp)
		return
	}
	// 校验成功，继续
	log.Info("时间:%s--请求来自:%s--用户ID:%d", time.Now().String(), req.Request.RemoteAddr, uid)
	req.SetAttribute("uid", uid)
	chain.ProcessFilter(req, resp)
}
