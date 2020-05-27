// Author: yann
// Date: 2019/9/26 上午10:39

package redisClient

import (
	"testing"
)

//不用测试,测试短信条数有限制,等集成测试的时候再用
func TestRedisPool(t *testing.T) {
	//config := &RedisConfig{"127.0.0.1:6379", 20}
	//Init(config)
	//rediss := RedisPool.Get()
	//
	////设置 并设置超时15秒
	//_, err := rediss.Do("SADD", "seen_video:123456", 10001,10002,10003,10004)
	//if err != nil{
	//	fmt.Println(err)
	//	return
	//}
	//
	////获取 转成string
	//reply, err := redis.Int64s(rediss.Do("SMEMBERS", "seen_video:123456"))
	//if err != nil{
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(reply)
}
