// Author: yann
// Date: 2019/9/26 上午10:30

package redisClient

import (
	"context"
	"errors"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"sync"
	chat "yann-chat"
)

type RedisConfig struct {
	Addr string
	Size int
}

var (
	redispool *redigo.Pool
	once      sync.Once
	ServerKey string
)

func GetClient() redigo.Conn {
	return redispool.Get()
}

//初始化参数
func (r *RedisConfig) initParam() {
	if r.Addr == "" {
		r.Addr = "127.0.0.1:6379"
	}
	if r.Size == 0 {
		r.Size = 128
	}
}

func (r *RedisConfig) Start(ctx context.Context, yannChat *chat.YannChat) error {
	once.Do(func() {
		r.initParam()
		redispool = redigo.NewPool(func() (conn redigo.Conn, e error) {
			conn, e = redigo.Dial("tcp", r.Addr)
			if e != nil {
				panic(errors.New("redis 连接池 初始化失败"))
			}
			return
		}, r.Size)
		logrus.Infof("redis %s 已启动", r.Addr)
	})
	return nil
}

func (r *RedisConfig) Stop(ctx context.Context) (err error) {
	err = redispool.Close()
	if err != nil {
		logrus.Errorf("redis 关闭异常: %s", err.Error())
	}
	logrus.Infof("redis 已关闭")
	return
}
