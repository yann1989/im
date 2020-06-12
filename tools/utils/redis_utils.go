// Author: yann
// Date: 2020/5/23 7:36 上午
// Desc:

package utils

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"reflect"
	"strconv"
	"time"
	"yann-chat/common"
	"yann-chat/tools/dao/redisClient"
)

type redisUtils struct{}

var RedisUtils = redisUtils{}

// 通过key获取string
func (redisUtils) GetStringByKey(key string) string {
	args := []interface{}{}
	args = append(args, "GET")
	args = append(args, key)
	cmd := redis.NewStringCmd(args...)
	var err error
	if redisClient.IsCluster {
		err = redisClient.ClusterClient.Process(cmd)
	} else {
		err = redisClient.Client.Process(cmd)
	}
	if err != nil {
		logrus.Errorf("redis GET error: %s", err.Error())
		return ""
	}
	result, err := cmd.Result()
	if err != nil {
		logrus.Errorf("redis GET error: %s", err.Error())
		return ""
	}
	return result
}

// 自增返回当前数字
func (redisUtils) IncrRedisByKeyAndGetValue(key string) int64 {
	args := []interface{}{}
	args = append(args, "INCR")
	args = append(args, key)
	cmd := redis.NewIntCmd(args...)
	var err error
	if redisClient.IsCluster {
		err = redisClient.ClusterClient.Process(cmd)
	} else {
		err = redisClient.Client.Process(cmd)
	}
	if err != nil {
		logrus.Errorf("redis GET error: %s", err.Error())
		return -1
	}
	num, err := cmd.Result()
	if err != nil {
		logrus.Errorf("redis INCR error: %s", err.Error())
		return -1
	}
	return num
}

// 设置string
func (redisUtils) SetStringByKey(key string, value interface{}, time int) error {
	args := []interface{}{}
	args = append(args, "SET")
	args = append(args, key)
	args = append(args, value)
	if time > 0 {
		args = append(args, "EX")
		args = append(args, time)
	}

	if redisClient.IsCluster {
		return redisClient.ClusterClient.Do(args...).Err()
	}
	return redisClient.Client.Do(args...).Err()

}

// 删除
func (redisUtils) DelByKey(key string) error {
	args := []interface{}{}
	args = append(args, "DEL")
	args = append(args, key)

	if redisClient.IsCluster {
		return redisClient.ClusterClient.Do(args...).Err()
	}
	return redisClient.Client.Do(args...).Err()
}

// redis SADD  set集合添加元素
func (redisUtils) Sadd(key string, args ...interface{}) error {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "SADD")
	tempArgs = append(tempArgs, key)
	tempArgs = append(tempArgs, args...)
	if redisClient.IsCluster {
		return redisClient.ClusterClient.Do(tempArgs...).Err()
	}
	return redisClient.Client.Do(tempArgs...).Err()
}

// redis Srem  set移除元素
func (redisUtils) Srem(key string, args ...interface{}) error {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "SREM")
	tempArgs = append(tempArgs, key)
	tempArgs = append(tempArgs, args...)
	if redisClient.IsCluster {
		return redisClient.ClusterClient.Do(tempArgs...).Err()
	}
	return redisClient.Client.Do(tempArgs...).Err()
}

// redis Smembers 获取set中所有元素
func (redisUtils) SmembersInt64(key interface{}) ([]int64, error) {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "SMEMBERS")
	tempArgs = append(tempArgs, key)
	cmd := redis.NewStringSliceCmd(tempArgs...)
	var err error
	if redisClient.IsCluster {
		err = redisClient.ClusterClient.Process(cmd)
	} else {
		err = redisClient.Client.Process(cmd)
	}
	if err != nil {
		return nil, err
	}
	//result, err := cmd.Result()
	array := []int64{}

	//for _, str := range result {
	//	int64Num, _ := strconv.ParseInt(str, 10, 64)
	//	if int64Num != 0 {
	//		array = append(array, int64Num)
	//	}
	//}
	return array, cmd.ScanSlice(&array)
}

// redis Smembers 获取set中所有元素
func (redisUtils) SmembersString(key interface{}) ([]string, error) {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "SMEMBERS")
	tempArgs = append(tempArgs, key)
	cmd := redis.NewStringSliceCmd(tempArgs...)
	var err error
	if redisClient.IsCluster {
		err = redisClient.ClusterClient.Process(cmd)
	} else {
		err = redisClient.Client.Process(cmd)
	}
	if err != nil {
		return nil, err
	}
	return cmd.Result()
}

// redis hash类型设置,设置单个字段
func (redisUtils) Hset(key string, field, value interface{}) error {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "HSET")
	tempArgs = append(tempArgs, key)
	tempArgs = append(tempArgs, field)
	tempArgs = append(tempArgs, value)
	if redisClient.IsCluster {
		return redisClient.ClusterClient.Do(tempArgs...).Err()
	}
	return redisClient.Client.Do(tempArgs...).Err()
}

// redis hash类型设置,批量设置字段
func (redisUtils) Hmset(key string, data interface{}) error {
	if redisClient.IsCluster {
		return redisClient.ClusterClient.HMSet(key, struct2Map(data)).Err()
	}
	return redisClient.Client.HMSet(key, struct2Map(data)).Err()
}

func struct2Map(i interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	elem := reflect.ValueOf(i).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return m
}

//
//// redis hash类型值获取
//func (redisUtils) HgetAll(key string) (map[string]interface{}, error) {
//	m := make(map[string]interface{})
//	value, err := redis.Strings(redisClient.GetClient().Do("HGETALL", key))
//	if err != nil {
//		return m, err
//	}
//	for k, v := range value {
//		if k%2 == 0 { // 如果是偶数,则是key,反之是value
//			m[v] = value[k+1]
//		}
//	}
//	return m, nil
//}
//
//// 获取hash里单个值
//func (redisUtils) Hget(key, field string) (string, error) {
//	value, err := redis.String(redisClient.GetClient().Do("HGET", key, field))
//	return value, err
//}

//***************************************************
//Description : 获取redis hash类型所有字段
//param :       key
//param :       对应结构体指针
//return :      错误信息
//***************************************************
func (redisUtils) HgetAll2Struct(key string, i interface{}) error {
	var cmd *redis.StringStringMapCmd = nil
	if redisClient.IsCluster {
		cmd = redisClient.ClusterClient.HGetAll(key)
	} else {
		cmd = redisClient.Client.HGetAll(key)
	}

	result, err := cmd.Result()
	if err != nil {
		return err
	}

	elem := reflect.ValueOf(i).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		if value, ok := result[relType.Field(i).Name]; ok {
			conversion, err := typeConversion(elem.Field(i).Type().Name(), value)
			if err != nil {
				continue
			}
			elem.Field(i).Set(conversion)
		}
	}
	return nil
}

func typeConversion(typeName, value string) (reflect.Value, error) {
	switch typeName {
	case "string":
		return reflect.ValueOf(value), nil
	case "time.Time":
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	case "Time":
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	case "int":
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	case "int8":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	case "int32":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int32(i)), err
	case "int64":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	case "float32":
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	case "float64":
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}
	return reflect.ValueOf(value), errors.New("未知的类型：" + typeName)
}

//
// redis zset 类型设置, 单个设置
func (redisUtils) ZaddSingle(key string, score int64, value string) error {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "ZADD")
	tempArgs = append(tempArgs, key)
	tempArgs = append(tempArgs, score)
	tempArgs = append(tempArgs, value)
	if redisClient.IsCluster {
		return redisClient.ClusterClient.Do(tempArgs...).Err()
	}
	return redisClient.Client.Do(tempArgs...).Err()
}

// redis zset 类型设置, 单个移除
func (redisUtils) ZremSingle(key string, value string) error {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "ZREM")
	tempArgs = append(tempArgs, key)
	tempArgs = append(tempArgs, value)
	if redisClient.IsCluster {
		return redisClient.ClusterClient.Do(tempArgs...).Err()
	}
	return redisClient.Client.Do(tempArgs...).Err()
}

// redis zset 类型设置, 获取score
func (redisUtils) Zscore(key string, value string) (int64, error) {
	args := []interface{}{}
	args = append(args, "ZSCORE")
	args = append(args, key)
	args = append(args, value)
	cmd := redis.NewStringCmd(args...)
	var err error
	if redisClient.IsCluster {
		err = redisClient.ClusterClient.Process(cmd)
	} else {
		err = redisClient.Client.Process(cmd)
	}
	if err != nil {
		logrus.Errorf("redis ZSCORE error: %s", err.Error())
		return 0, err
	}

	return strconv.ParseInt(cmd.Val(), 10, 64)
}

// redis 分页获取zset 列表
func (redisUtils) ZRANGEBYSCORE(key string, startScore, endScore int64, page, rows int) ([]string, error) {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "ZRANGEBYSCORE")
	tempArgs = append(tempArgs, key)
	tempArgs = append(tempArgs, startScore)
	tempArgs = append(tempArgs, endScore)
	tempArgs = append(tempArgs, "LIMIT")
	tempArgs = append(tempArgs, page)
	tempArgs = append(tempArgs, rows)
	cmd := redis.NewStringSliceCmd(tempArgs...)
	var err error
	if redisClient.IsCluster {
		err = redisClient.ClusterClient.Process(cmd)
	} else {
		err = redisClient.Client.Process(cmd)
	}
	if err != nil {
		return nil, err
	}
	return cmd.Result()
}

// redis 分页获取zset 列表,job用的
// -inf = 小于当前时间
func (redisUtils) ZRANGEBYSCOREBYJOB(key string, endScore int64, page, rows int) ([]string, error) {
	var cmd *redis.StringSliceCmd = nil
	if redisClient.IsCluster {
		cmd = redisClient.ClusterClient.ZRangeByScore(key, &redis.ZRangeBy{Max: strconv.FormatInt(endScore, 10), Offset: int64(page * rows), Count: int64(rows)})
	} else {
		cmd = redisClient.Client.ZRangeByScore(key, &redis.ZRangeBy{Max: strconv.FormatInt(endScore, 10), Offset: int64(page * rows), Count: int64(rows)})
	}
	return cmd.Result()
}

// redis 获取zset
func (redisUtils) ZRANGE(key string) (map[string]string, error) {
	tempArgs := []interface{}{}
	tempArgs = append(tempArgs, "ZRANGE")
	tempArgs = append(tempArgs, key)
	tempArgs = append(tempArgs, 0)
	tempArgs = append(tempArgs, 1)
	tempArgs = append(tempArgs, "WITHSCORES")
	cmd := redis.NewStringStringMapCmd(tempArgs...)
	var err error
	if redisClient.IsCluster {
		err = redisClient.ClusterClient.Process(cmd)
	} else {
		err = redisClient.Client.Process(cmd)
	}
	if err != nil {
		return nil, err
	}
	return cmd.Result()
}

func (redisUtils) BuildKey(prefix common.RedisKey, key interface{}) string {
	redisKey := string(prefix)
	redisKey += cast.ToString(key)
	return redisKey
}
