package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	chat "yann-chat"
	"yann-chat/common"
	"yann-chat/tools/conf"
	"yann-chat/tools/ip"
	"yann-chat/tools/log"
)

// package-private vars
var (
	// 配置信息总控, main package 私有, 不要对外暴露, 所有的配置加载都在main
	config *conf.Config

	// 配置文件路径, 通过命令行flag进行加载
	configFile string

	// 所有子命令的root command
	rootCmd = &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "the " + filepath.Base(os.Args[0]) + "command line interface",
		Long: "This is the " + filepath.Base(os.Args[0]) +
			` command line interface.For details on the dao,
please contact the relevant technical staff.`,
		Version: version(nil),
	}

	// 子命令-启动入口命令
	startCmd = &cobra.Command{
		Use:   "start [command]",
		Short: "start the dao",
		Run:   entryPoint,
	}

	// 进程中的单例, 负责注册管理所有的子服务, 比如web层, rpc层, 等等.
	yannChat = chat.NewYannChat()
)

// 程序入口点, 在此处添加所有启动所需要的逻辑
func entryPoint(_ *cobra.Command, _ []string) {

	var err error

	// configuration
	loadConfig()

	// 注册雪花id
	if err := yannChat.Register(func(ctx *context.Context) (service chat.Service, err error) {
		return config.Snowflake, nil
	}); err != nil {
		panic(err)
	}
	// 注册redis
	if err := yannChat.Register(func(ctx *context.Context) (service chat.Service, err error) {
		return config.Dao.Redis, nil
	}); err != nil {
		panic(err)
	}
	// 注册mq
	if err := yannChat.Register(func(ctx *context.Context) (service chat.Service, err error) {
		return config.MQ, nil
	}); err != nil {
		panic(err)
	}
	// 构建连接管理工具
	if err := yannChat.Register(func(ctx *context.Context) (service chat.Service, err error) {
		return config.Manager, nil
	}); err != nil {
		panic(err)
	}
	// 注册web
	if err := yannChat.Register(func(ctx *context.Context) (service chat.Service, err error) {
		return NewHttpServer(config), nil
	}); err != nil {
		panic(err)
	}

	// 构建Taurus
	value := context.WithValue(context.Background(), "conf", config)
	if err = yannChat.Build(value); err != nil {
		log.Error("Building Taurus failed, please check the startup settings. err: %s", err.Error())
		panic(err)
	}

	// 启动
	if err = yannChat.Start(value); err != nil {
		log.Error("Launching Taurus failed, please check the startup settings. err: %s", err.Error())
		panic(err)
	}

	// 启动时间日志记录
	log.Info("chat-service starts at %s", config.Web.Address)
	log.Info("启动时间=>[%s]", time.Now().Format("2006-01-02 15:04:05"))

	// Service Shutdown
	shutdown()
}

func Execute() error {
	return rootCmd.Execute()
}

/* Private helper functions */

// 加载配置
func loadConfig() {

	// Get mini_conf file location
	if configFile == "" {
		fmt.Println("Please use ./" + filepath.Base(os.Args[0]) + " start --config config.toml")
		os.Exit(-1)
	}

	// load mini_conf info from mini_conf file to mini_conf
	var err error
	if config, err = conf.Init(configFile); err != nil {
		fmt.Println("Failed to parse the configuration file, please check the configuration file")
		os.Exit(-1)
	}

	// init log
	config.Log.ServiceAddress = ip.InternalIP()
	log.Init(config.Log)
}

// 服务停止
func shutdown() {
	var err error
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("chat-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			common.IsClose = true
			time.Sleep(time.Second * 2)
			if err = yannChat.Stop(context.Background()); err != nil {
				log.Error("服务退出异常: %v", err)
			}
			log.Info("服务结束运行")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
