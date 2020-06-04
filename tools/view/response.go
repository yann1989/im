// Author       yann
// Time         2019-11-17 10:51
// File Desc    前端响应

package view

import (
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"time"
)

// 响应前端的固定格式
type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

//服务异常 统一返回服务忙
func Response500() (response *Response) {
	response = new(Response)
	response.Code = CODE_SERVER_ERR
	response.Msg = ERR_SERVER
	response.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	response.Data = struct{}{}
	return
}

//token解析错误
func Response401(request *restful.Request) (response *Response) {
	response = new(Response)
	response.Code = CODE_TOKEN_ERR
	response.Msg = ERR_TOKEN
	response.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	response.Data = struct{}{}
	return
}

// 响应前端数据
func (r *Response) ReturnResult(request *restful.Request, response *restful.Response) {
	//设置跨域
	response.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	response.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Lang, Authorization")
	if err := response.WriteAsJson(r); err != nil {
		logrus.Errorf("%s返回结果失败", request.Request.URL.Path)
	}
}
