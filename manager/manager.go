// Author: yann
// Date: 2020/5/23 10:13 上午
// Desc:

package manager

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"
	"github.com/sirupsen/logrus"
	"sync"
	chat "yann-chat"
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
		logrus.Errorf(err)
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
			logrus.Errorf("init poller faild:", err.Error())
			panic("init poller faild:" + err.Error())
		}
		manager.gopool = NewPool(GOROUTING_MAX_LEN, TASK_MAX_LEN, GOROUTING_INIT_LEN)
		go m.startConsume()
		logrus.Infof("connect manager[max_connect:%d] 初始化完成", manager.MaxConnect)
	})
	return nil
}

//停止
func (m *ConnectManager) Stop(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("释放所有websocket 未知错误: %S", err)
		}
	}()
	for key, _ := range manager.nodes {
		manager.Remove(key)
	}
	logrus.Infof("websocket所有连接 已关闭")
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
	//失败会panic, 请在调用方法前recover
	fd := netpoll.Must(netpoll.HandleRead(conn.UnderlyingConn()))
	node.fd = fd
	if err := manager.poller.Start(fd, func(ev netpoll.Event) {
		if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("eof 释放资源panic: %s", err)
				}
			}()
			m.Remove(id) //闭包
		}

		logrus.Infof("读事件触发")
		manager.gopool.Schedule(func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("读事件执行失败")
				}
			}()
			Receiver(node) //闭包

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
//Description : 从内核树中移除, 关闭连接, 删除节点
//param :       用户id
//***************************************************
func (m *ConnectManager) Remove(id int64) {
	manager.rwlocker.Lock()
	defer manager.rwlocker.Unlock()
	if node, has := manager.nodes[id]; has {
		err := manager.poller.Stop(node.fd)
		if err != nil {
			logrus.Errorf("关闭连接失败: %S", err)
		}
		delete(manager.nodes, id)
		manager.currentConnect--
	}

}
