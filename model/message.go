// Author: yann
// Date: 2020/5/23 12:06 下午
// Desc:

package model

type Message struct {
	MsgID       int64  `json:"msg_id,omitempty"`       //唯一id
	SessionID   int64  `json:"session_id,omitempty"`   //私聊的时候必会话id
	FromID      int64  `json:"from_id,omitempty"`      //谁发的
	ToID        int64  `json:"to_id,omitempty"`        //发给谁
	MsgType     int    `json:"msg_type"`               //消息类型 1=群聊 2=私聊
	Content     string `json:"content,omitempty"`      //文本内容
	ContentType int    `json:"content_type,omitempty"` //内容类型
	ResourceUrl string `json:"resource_url,omitempty"` //图片,视频,音频 url
	CreateTime  int64  `json:"create_time,omitempty"`  //发送时间
	Memo        string `json:"memo,omitempty"`         //json 备用字段
}
