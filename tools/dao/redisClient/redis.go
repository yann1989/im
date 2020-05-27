// Author: yann
// Date: 2019/9/26 上午10:30

package redisClient

import (
	"context"
	"errors"
	redigo "github.com/gomodule/redigo/redis"
	"sync"
	chat "yann-chat"
	"yann-chat/common"
	"yann-chat/tools/log"
)

type RedisConfig struct {
	Addr      string
	Size      int
	ServerKey string
}

var (
	Conn      redigo.Conn
	Conn2     redigo.Conn
	once      sync.Once
	ServerKey string
)

//初始化参数
func (r *RedisConfig) initParam() {
	if r.Addr == "" {
		r.Addr = "127.0.0.1:6379"
	}
	//if r.Size == 0 {
	//	r.Size = 100
	//}
	if r.ServerKey == "" {
		panic("请初始化ServerKey")
	}
	r.ServerKey = string(common.REDIS_KEY_MESSAGE_QUEUE_PREFIX) + r.ServerKey
	ServerKey = r.ServerKey
}

func (r *RedisConfig) Start(ctx context.Context, yannChat *chat.YannChat) error {
	once.Do(func() {
		var err error
		r.initParam()
		Conn, err = redigo.Dial("tcp", r.Addr)
		if err != nil {
			panic(errors.New("redis 初始化失败"))
		}
		log.Info("redis1 %s 已启动", r.Addr)
		if _, err := Conn.Do("SADD", common.REDIS_KEY_IM_SERVER, r.ServerKey); err != nil {
			panic("redis sadd im:service 写入im节点:" + r.ServerKey + "失败")
		}
		Conn2, err = redigo.Dial("tcp", r.Addr)
		if err != nil {
			panic(errors.New("redis 初始化失败"))
		}
		log.Info("redis2 %s 已启动", r.Addr)
		log.Info("redis sadd im:service 节点:%s 已写入", r.ServerKey)
		return
	})
	return nil
}

func (r *RedisConfig) Stop(ctx context.Context) (err error) {
	if _, err := Conn.Do("SREM", common.REDIS_KEY_IM_SERVER, r.ServerKey); err != nil {
		log.Error("redis SREM im:service 移除 im节点:" + r.ServerKey + "失败")
	}
	err = Conn.Close()
	if err != nil {
		log.Error("redis 关闭异常: %s", err.Error())
	}
	log.Info("redis 已关闭")
	return
}
