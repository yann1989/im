// Author       kevin
// Time         2019-08-08 20:16
// File Desc    配置

package conf

import (
	"github.com/BurntSushi/toml"
	"yann-chat/manager"
	"yann-chat/tools/dao/mysqlClient"
	"yann-chat/tools/dao/redisClient"
	"yann-chat/tools/log"
	"yann-chat/tools/mq"
	"yann-chat/tools/snowflake"
	"yann-chat/tools/web"
)

// Config 服务配置总控-对应配置文件所有的配置信息
type Config struct {
	// 日志
	Log *log.Config `toml:"log"`
	// 数据层配置
	Dao *DaoConfig `toml:"dao"`
	// web
	Web *web.Config `toml:"web"`
	// 雪花
	Snowflake *snowflake.SnowConfig `toml:"snowflake"`
	// 客户端管理
	Manager *manager.ConnectManager `toml:"manager"`
	//mq
	MQ *mq.Config `toml:"mq"`
}

//数据层配置
type DaoConfig struct {
	Mysql *mysqlClient.MysqlConfig
	Redis *redisClient.RedisConfig
}

// 加载配置文件信息
// [参数]
// filePath, 配置文件位置
// [返回值]
// 配置总控信息, 错误信息
func Init(filePath string) (conf *Config, err error) {
	conf = new(Config)
	_, err = toml.DecodeFile(filePath, &conf)
	if err != nil {
		return
	}
	return
}
