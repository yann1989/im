// Author: yann
// Date: 2020/5/23 7:36 上午
// Desc:

package utils

import (
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/cast"
	"strings"
	"yann-chat/common"
	"yann-chat/tools/dao/redisClient"
)

import (
	redigo "github.com/gomodule/redigo/redis"
)

type redisUtils struct{}

var RedisUtils = redisUtils{}

// 通过key获取string
func (redisUtils) GetStringByKey(key string) string {
	// 连接redis
	value, err := redis.String(redisClient.Conn.Do("GET", key))
	if err != nil {
		return ""
	}
	return value
}

// 自增返回当前数字
func (redisUtils) IncrRedisByKeyAndGetValue(key string) (int, error) {
	return redis.Int(redisClient.Conn.Do("INCR", key))
}

// 设置string
func (redisUtils) SetStringByKey(key string, value interface{}, time int) error {
	// 连接redis
	var err error
	if time == 0 { // 设置不过期key
		_, err = redisClient.Conn.Do("SET", key, value)
	} else {
		_, err = redisClient.Conn.Do("SET", key, value, "EX", time)
	}
	return err
}

// 删除
func (redisUtils) DelByKey(key string) error {
	_, err := redisClient.Conn.Do("DEL", key)
	return err
}

// redis SADD  set集合添加元素
func (redisUtils) Sadd(key string, args ...interface{}) error {
	arr := []interface{}{key}
	for _, value := range args {
		arr = append(arr, value)
	}
	_, err := redisClient.Conn.Do("SADD", arr...)
	return err
}

// redis Srem  set移除元素
func (redisUtils) Srem(key string, args ...interface{}) error {
	arr := []interface{}{key}
	for _, value := range args {
		arr = append(arr, value)
	}
	_, err := redisClient.Conn.Do("SREM", arr...)
	return err
}

// redis Smembers 获取set中所有元素
func (redisUtils) SmembersInt64(key interface{}) []int64 {
	list, err := redis.Int64s(redisClient.Conn.Do("SMEMBERS", key))
	if err != nil {
		return nil
	}
	return list
}

// redis Smembers 获取set中所有元素
func (redisUtils) SmembersString(key interface{}) ([]string, error) {
	return redis.Strings(redisClient.Conn2.Do("SMEMBERS", key))
}

// redis hash类型设置,设置单个字段
func (redisUtils) Hset(key string, field, value interface{}) error {
	_, err := redisClient.Conn2.Do("HSET", key, field, value)
	return err
}

// redis hash类型设置,批量设置字段
func (redisUtils) Hmset(key string, data interface{}) error {
	_, err := redisClient.Conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(data)...)
	return err
}

// redis hash类型值获取
func (redisUtils) HgetAll(key string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	value, err := redis.Strings(redisClient.Conn.Do("HGETALL", key))
	if err != nil {
		return m, err
	}
	for k, v := range value {
		if k%2 == 0 { // 如果是偶数,则是key,反之是value
			m[v] = value[k+1]
		}
	}
	return m, nil
}

// 获取hash里单个值
func (redisUtils) Hget(key, field string) (string, error) {
	value, err := redis.String(redisClient.Conn.Do("HGET", key, field))
	return value, err
}

//***************************************************
//Description : 获取redis hash类型所有字段
//param :       key
//param :       对应结构体指针
//return :      错误信息
//***************************************************
func (redisUtils) HgetAll2Struct(key string, i interface{}) error {
	value, err := redis.Values(redisClient.Conn.Do("HGETALL", key))
	if err != nil {
		return err
	}
	return redis.ScanStruct(value, i)
}

// redis zset 类型设置, 单个设置
func (redisUtils) ZaddSingle(key string, score int64, value string) error {
	_, err := redisClient.Conn2.Do("ZADD", key, score, value)
	return err
}

// redis zset 类型设置, 单个移除
func (redisUtils) ZremSingle(key string, value string) error {
	_, err := redisClient.Conn.Do("ZREM", key, value)
	return err
}

// redis zset 类型设置, 获取score
func (redisUtils) Zscore(key string, value string) (string, error) {
	score, err := redisClient.Conn.Do("ZSCORE", key, value)
	if score == nil {
		return "", err
	}
	return string(score.([]byte)), err
}

// redis 分页获取zset 列表
func (redisUtils) ZRANGEBYSCORE(key string, startScore, endScore int64, page, rows int) ([]string, error) {
	var value []string
	value, err := redis.Strings(redisClient.Conn.Do("ZRANGEBYSCORE", key, startScore, endScore, "LIMIT", page, rows))
	return value, err
}

// redis 分页获取zset 列表,job用的
// -inf = 小于当前时间
func (redisUtils) ZRANGEBYSCOREBYJOB(key string, endScore int64, page, rows int) ([]string, error) {
	var value []string
	value, err := redis.Strings(redisClient.Conn.Do("ZRANGEBYSCORE", key, "-inf", endScore, "LIMIT", page, rows))
	return value, err
}

// redis 获取zset
func (redisUtils) ZRANGE(key string) (map[string]string, error) {
	value, err := redis.StringMap(redisClient.Conn.Do("ZRANGE", key, 0, 1, "WITHSCORES"))
	return value, err
}

// redis 获取zset 交集(只求2个key)
func (redisUtils) ZINTERSTORE(key1, key2, newKey string) int {
	//v, err := redisCilent.Conn.Do()("ZINTERSTORE", newKey, 2, key1, key2)
	i, _ := redisClient.Conn.Do("ZINTERSTORE", newKey, 2, key1, key2)
	return cast.ToInt(i)
}

// redis zset 集合成员数
func (redisUtils) ZCARD(key string) (int, error) {
	s, err := redis.Int(redisClient.Conn.Do("ZCARD", key))
	return s, err
}

//读取redis中自己的消息队列
func (redisUtils) LpopMessage() []byte {
	msg, _ := redis.Bytes(redisClient.Conn.Do("LPOP", redisClient.ServerKey))
	return msg
}

// 存储单个list
func (redisUtils) LPUSH(key, value string) error {
	_, err := redisClient.Conn.Do("LPUSH", key, value)
	return err
}

// 移除单个元素
func (redisUtils) BLPOP(key string, conn redigo.Conn) error {
	_, err := redisClient.Conn.Do("BLPOP", key, 3)
	return err
}

// 读取list值
func (redisUtils) LRANGE(key string) (string, error) {
	count, _ := redisClient.Conn.Do("LLEN", key)
	value, _ := redis.Values(redisClient.Conn.Do("LRANGE", key, 0, count))
	var data []string
	if err := redis.ScanSlice(value, &data); err != nil {
		return "", err
	}
	// 如果为空
	if data == nil || len(data) == 0 {
		return "", nil
	}
	// 拼接成json格式
	return "[" + strings.Join(data, ",") + "]", nil
}

//可以一次push多条数据
func (redisUtils) RPUSH(key string, args ...interface{}) error {
	arr := []interface{}{key}
	for _, value := range args {
		arr = append(arr, value)
	}
	_, err := redisClient.Conn.Do("RPUSH", arr...)
	return err
}

//获取list长度
func (redisUtils) LLEN(key string) (int, error) {
	return redis.Int(redisClient.Conn.Do("LLEN", key))
}

func (redisUtils) BuildKey(prefix common.RedisKey, key interface{}) string {
	redisKey := string(prefix)
	redisKey += cast.ToString(key)
	return redisKey
}
