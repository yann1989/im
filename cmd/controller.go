// Author       kevin
// Time         2019-08-08 20:16
// File Desc    controller层主入口

package main

import (
	"crypto/ecdsa"
	"io/ioutil"
	"net/http"
	"strings"
	chat "yann-chat"
	"yann-chat/tools/conf"
	"yann-chat/tools/ip"
	"yann-chat/tools/jwt"
	"yann-chat/tools/log"

	"github.com/emicklei/go-restful"
	"golang.org/x/net/context"
)

type HttpServer struct {
	// protocol: http & https
	proto string

	// go-restful框架的container: 负责路由
	*restful.Container

	// go native http service
	server *http.Server

	// 公私钥文件位置
	certPath, keyPath string

	// 私钥, 公钥通过私钥获取 - privateKey.PublicKey
	privateKey *ecdsa.PrivateKey
}

// 创建Web层
// [参数]
// conf: 配置信息
// [返回值]
// httpServer
func NewHttpServer(config *conf.Config) (httpServer *HttpServer) {

	// 默认监听绑定内网IPv4地址的8901端口
	if config.Web.Address == "" {
		config.Web.Address = ip.InternalIP() + ":8901"
	}

	if config.Web.Proto == "" {
		config.Web.Proto = "https"
	}

	// 私钥
	// private key
	privateKeyBytes, err := ioutil.ReadFile(config.Web.KeyPath)
	priKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Error("注册web服务失败: %v", err)
		panic(err)
	}

	// 创建一个httpServer
	server := &http.Server{Addr: config.Web.Address, Handler: http.DefaultServeMux}

	httpServer = &HttpServer{
		proto:      config.Web.Proto,
		Container:  restful.NewContainer(),
		server:     server,
		certPath:   config.Web.CertPath,
		keyPath:    config.Web.KeyPath,
		privateKey: priKey,
	}

	ws := new(restful.WebService)
	httpServer.initRoutes(ws)
	httpServer.Container.Add(ws)

	// CORS
	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept", "Lang", "Authorization", "Device-Uuid"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      httpServer.Container}
	httpServer.Container.Filter(cors.Filter)
	// Add container filter to respond to OPTIONS
	httpServer.Container.Filter(httpServer.Container.OPTIONSFilter)

	return
}

// 启动http service
func (h *HttpServer) Start(ctx context.Context, yannChat *chat.YannChat) error {

	switch strings.ToLower(h.proto) {
	case "https":
		go func() {
			// returns ErrServerClosed on graceful shutdown(i.e. call service.Shutdown())
			if err := h.server.ListenAndServeTLS(h.certPath, h.keyPath); err != http.ErrServerClosed {
				log.Error("ListenAndServeTLS error: %v", err)
				panic("web层启动异常")
			}
		}()
	case "http":
		go func() {
			// returns ErrServerClosed on graceful shutdown(i.e. call service.Shutdown())
			if err := http.ListenAndServe(h.server.Addr, h.Container); err != http.ErrServerClosed {
				log.Error("ListenAndServe error: %v", err)
				panic("web层启动异常")
			}
		}()
	default:
		log.Error("无法识别协议: %s", h.proto)
		panic("无法识别协议, 请检查配置")
	}

	return nil
}

// 停止http service
func (h *HttpServer) Stop(ctx context.Context) (err error) {
	// 关闭http service
	if err = h.server.Shutdown(ctx); err != nil {
		log.Error("web层关闭异常: %v", err)
	}
	log.Info("web 服务已关闭")
	return
}
