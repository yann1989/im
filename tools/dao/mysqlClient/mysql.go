// Author: yann
// Date: 2019/9/21 上午10:51

package mysqlClient

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
	chat "yann-chat"
	"yann-chat/tools/log"
)

type MysqlConfig struct {
	DSNAPP             string
	DSNADMIN           string
	Debug              bool
	MaxIdleConns       int
	MaxOpenConns       int
	MaxConnMaxLifetime int64
}

var (
	DbApp   *gorm.DB
	DbAdmin *gorm.DB
	once    sync.Once
)

func (m *MysqlConfig) Start(ctx context.Context, yannChat *chat.YannChat) error {
	once.Do(func() {
		if m.DSNAPP == "" && m.DSNADMIN == "" {
			panic("请配置初始化数据库参数")
		}
		var err error
		if m.DSNAPP != "" {
			DbApp, err = gorm.Open("mysql", m.DSNAPP)
			if err != nil {
				panic(err)
			}
			DbApp.DB().SetMaxIdleConns(m.MaxIdleConns)
			DbApp.DB().SetMaxOpenConns(m.MaxOpenConns)
			DbApp.DB().SetConnMaxLifetime(time.Duration(m.MaxConnMaxLifetime) * time.Second)
			if m.Debug {
				DbApp = DbApp.Debug()
			}
		}
		log.Info("mysql %s已启动", m.DSNAPP)

		if m.DSNADMIN != "" {
			DbAdmin, err = gorm.Open("mysql", m.DSNADMIN)
			if err != nil {
				panic(err)
			}
			DbAdmin.DB().SetMaxIdleConns(m.MaxIdleConns)
			DbAdmin.DB().SetMaxOpenConns(m.MaxOpenConns)
			DbAdmin.DB().SetConnMaxLifetime(time.Duration(m.MaxConnMaxLifetime) * time.Second)
			if m.Debug {
				DbAdmin = DbAdmin.Debug()
			}
		}
		log.Info("mysql %s已启动", m.DSNADMIN)
	})
	log.Info("mysql 数据库已启动")
	return nil
}

func (m *MysqlConfig) Stop(ctx context.Context) (err error) {
	err = DbApp.Close()
	if err != nil {
		log.Error("mysql db_app 关闭异常: %s", err.Error())
		return
	}
	err = DbApp.Close()
	if err != nil {
		log.Error("mysql db_admin 关闭异常: %s", err.Error())
		return
	}
	log.Info("mysql 数据库已关闭")
	return
}
