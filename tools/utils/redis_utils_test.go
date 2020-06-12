// Author: yann
// Date: 2020/6/11 12:57 下午
// Desc:

package utils

import (
	"context"
	"fmt"
	"testing"
	chat "yann-chat"
	"yann-chat/tools/dao/redisClient"
)

var (
	ctx  = context.Background()
	node = chat.NewYannChat()
)

func InitRedis() {
	config := &redisClient.RedisConfig{Addr: "127.0.0.1:6379", IsCluster: false}
	err := config.Start(ctx, node)
	if err != nil {
		panic(fmt.Sprintf("redis初始化失败:%s", err.Error()))
	}
}

func TestSetStringByKey(t *testing.T) {
	InitRedis()
	err := RedisUtils.SetStringByKey("go_redis_set", "test1", 5)
	if err != nil {
		panic(fmt.Sprintf("redis SET 失败 :%s", err.Error()))
	}
}

func TestGetStringByKey(t *testing.T) {
	InitRedis()
	str := RedisUtils.GetStringByKey("go_redis_set")
	fmt.Println(str)
}

func TestIncr(t *testing.T) {
	InitRedis()
	num := RedisUtils.IncrRedisByKeyAndGetValue("go_redis_incr")
	fmt.Println(num)
}

func TestDelByKey(t *testing.T) {
	InitRedis()
	err := RedisUtils.DelByKey("go_redis_set")
	if err != nil {
		panic(err)
	}
}

func TestSadd(t *testing.T) {
	InitRedis()
	err := RedisUtils.Sadd("go_redis_sadd", 1, 2, 3, 4, 5)
	if err != nil {
		panic(err)
	}
}

func TestSrem(t *testing.T) {
	InitRedis()
	err := RedisUtils.Srem("go_redis_sadd", "EX", 1)
	if err != nil {
		panic(err)
	}
}

func TestSmembersInt64(t *testing.T) {
	InitRedis()
	smembersString, err := RedisUtils.SmembersInt64("go_redis_sadd")
	if err != nil {
		panic(err)
	}
	fmt.Println(smembersString)
}

func TestSmembersString(t *testing.T) {
	InitRedis()
	smembersInt64, err := RedisUtils.SmembersString("go_redis_sadd")
	if err != nil {
		panic(err)
	}
	fmt.Println(smembersInt64)
}

type User struct {
	Name string
	Age  int
}

func TestHmset(t *testing.T) {
	InitRedis()
	err := RedisUtils.Hmset("go_redis_hset", &User{Name: "yann", Age: 18})
	if err != nil {
		panic(err)
	}
}

func TestHgetAll2Struct(t *testing.T) {
	InitRedis()
	user := &User{}
	err := RedisUtils.HgetAll2Struct("go_redis_hset", user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}

func TestZaddSingle(t *testing.T) {
	InitRedis()
	err := RedisUtils.ZaddSingle("go_redis_zset", 100, "aaaa")
	if err != nil {
		panic(err)
	}
}

func TestZscore(t *testing.T) {
	InitRedis()
	zscore, err := RedisUtils.Zscore("go_redis_zset", "aaaa")
	if err != nil {
		panic(err)
	}
	fmt.Println(zscore)
}

func TestZRANGEBYSCORE(t *testing.T) {
	InitRedis()
	zscore, err := RedisUtils.ZRANGEBYSCORE("go_redis_zset", 1, 101, 0, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(zscore)
}

func TestZRANGEBYSCOREBYJOB(t *testing.T) {
	InitRedis()
	zscore, err := RedisUtils.ZRANGEBYSCOREBYJOB("go_redis_zset", 100, 0, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(zscore)
}

func TestZRANGE(t *testing.T) {
	InitRedis()
	zscore, err := RedisUtils.ZRANGE("go_redis_zset")
	if err != nil {
		panic(err)
	}
	fmt.Println(zscore)
}
