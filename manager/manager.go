// Author: yann
// Date: 2020/5/23 10:13 上午
// Desc:

package manager

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"
	"sync"
	chat "yann-chat"
	"yann-chat/tools/log"
)

type ConnectManager struct {
	MaxConnect     int
	nodes          map[int64]*Node
	currentConnect int
	rwlocker       sync.RWMutex
	poller         netpoll.Poller
	gopool         *Pool
}

var (
	manager *ConnectManager //所有连接的管理者
	once    sync.Once
)

func (m *ConnectManager) verifyParam() {
	if m.MaxConnect <= 0 {
		err := fmt.Sprintf("客户端管理参数错误 MaxConnect = %d, 请检查配置", m.MaxConnect)
		log.Error(err)
		panic(err)
	}
}

//初始化
func (m *ConnectManager) Start(ctx context.Context, yannChat *chat.YannChat) error {
	once.Do(func() {
		m.verifyParam()
		manager = new(ConnectManager)
		manager.MaxConnect = m.MaxConnect
		manager.nodes = make(map[int64]*Node)
		manager.rwlocker = sync.RWMutex{}
		manager.currentConnect = 0
		var err error
		manager.poller, err = netpoll.New(nil)
		if err != nil {
			log.Error("init poller faild:", err.Error())
			panic("init poller faild:" + err.Error())
		}
		manager.gopool = NewPool(1024, 4096, 1)
		go m.startConsume()
		log.Info("connect manager[max_connect:%d] 初始化完成", manager.MaxConnect)
	})
	return nil
}

//停止
func (m *ConnectManager) Stop(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("eof 释放资源失败")
		}
	}()
	for key, _ := range manager.nodes {
		manager.Remove(key)
	}
	log.Info("websocket所有连接 已关闭")
	return nil
}

//请求到达判断连接数量不超过最大数量, 就可以尝试连接.
func (m *ConnectManager) isMax() bool {
	manager.rwlocker.RLock()
	defer manager.rwlocker.RUnlock()
	return manager.MaxConnect <= manager.currentConnect
}

//***************************************************
//Description : 将用户连接后的节点加入管理中
//param :       用户id
//param :       连接节点
//return :      成功true  失败 false(需要释放node节点)
//***************************************************
func (m *ConnectManager) Add(id int64, conn *websocket.Conn) *Node {
	//如果有未释放的资源则释放资源
	m.ifConnectedThanRemove(id)

	//再次判断当前连接数是否已满
	if manager.currentConnect >= manager.MaxConnect {
		return nil
	}
	node := new(Node)
	node.ClientId = id
	node.Conn = conn
	fd := netpoll.Must(netpoll.HandleRead(conn.UnderlyingConn()))
	node.fd = fd
	if err := manager.poller.Start(fd, func(ev netpoll.Event) {
		if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			defer func() {
				if err := recover(); err != nil {
					log.Error("eof 释放资源失败: %s", err)
				}
			}()
			log.Info("对端关闭事件触发")
			m.Remove(id)
			log.Info("websocket 被关闭")
		}

		log.Info("读事件触发")
		manager.gopool.Schedule(func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("读事件执行失败")
				}
			}()
			Receiver(node)

		})

	}); err != nil {
		return nil
	}
	manager.rwlocker.Lock()
	defer manager.rwlocker.Unlock()
	manager.nodes[id] = node
	manager.currentConnect++
	return node
}

//***************************************************
//Description : 如果已连接则删除之前的连接
//param :       用户id
//***************************************************
func (m *ConnectManager) ifConnectedThanRemove(id int64) {
	m.Remove(id)
}

//***************************************************
//Description : 删除节点,使用前请先加互斥锁
//param :       用户id
//***************************************************
func (m *ConnectManager) Remove(id int64) {
	manager.rwlocker.Lock()
	defer manager.rwlocker.Unlock()
	if node, has := manager.nodes[id]; has {
		err := manager.poller.Stop(node.fd)
		if err != nil {
			log.Error("重复登录 关闭之前连接失败")
		}
		delete(manager.nodes, id)
		manager.currentConnect--
	}

}
