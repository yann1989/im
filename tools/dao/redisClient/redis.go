// Author: yann
// Date: 2019/9/26 上午10:30

package redisClient

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"strings"
	chat "yann-chat"
)

type RedisConfig struct {
	Addr      string
	IsCluster bool
}

const REDIS_SUB_KEY = "chicha:chat:sub"

var (
	Client        *redis.Client        = nil
	ClusterClient *redis.ClusterClient = nil
	IsCluster     bool
	Ch            <-chan *redis.Message
)

func (r *RedisConfig) Start(ctx context.Context, yannChat *chat.YannChat) error {
	var err error
	if r.IsCluster {
		ClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: strings.Split(r.Addr, ","),
		})
		_, err = ClusterClient.Ping().Result()
		if err != nil {
			panic(errors.New("redis集群 初始化失败:" + err.Error()))
		}
		IsCluster = true
	} else {
		Client = redis.NewClient(&redis.Options{
			Addr: r.Addr, // use default Addr
		})
		_, err = Client.Ping().Result()
		if err != nil {
			panic(errors.New("redis 初始化失败:" + err.Error()))
		}
	}
	subscription()
	logrus.Infof("redis %s 初始化成功", r.Addr)
	return nil
}

func (r *RedisConfig) Stop(ctx context.Context) (err error) {
	if IsCluster {
		err = ClusterClient.Close()
		if err != nil {
			logrus.Errorf("redis集群 关闭异常: %s", err.Error())
		}
		return err
	}

	err = Client.Close()
	if err != nil {
		logrus.Errorf("redis 关闭异常: %s", err.Error())
	}
	return err

}

//***************************************************
//Description : 订阅, 用于消息广播
//***************************************************
func subscription() {
	var pubSub *redis.PubSub = nil
	if IsCluster {
		pubSub = ClusterClient.Subscribe(REDIS_SUB_KEY)
	} else {
		pubSub = Client.Subscribe(REDIS_SUB_KEY)
	}
	_, err := pubSub.Receive()
	if err != nil {
		panic(errors.New("redis Subscribe faild:" + err.Error()))
	}
	Ch = pubSub.Channel()
}
