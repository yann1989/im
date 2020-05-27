// Author: yann
// Date: 2020/5/23 3:28 下午
// Desc:

package manager

//模拟心跳 超时关闭连接
//func HeartBeating(node *Node) {
//	for {
//		select {
//		case <-node.HeartCh:
//			log.Info("收到心跳")
//		//超时1分钟没有数据可读, 则删除资源
//		case <-time.After(time.Second * 3):
//			defer func() {
//				if err := recover(); err != nil {
//					fmt.Println("eof 释放资源失败")
//				}
//			}()
//			manager.rwlocker.Lock()
//			defer manager.rwlocker.Unlock()
//			err := manager.poller.Stop(node.fd)
//			if err != nil {
//				log.Error("stop关闭失败: %s", err.Error())
//			}
//			err = node.fd.Close()
//			if err != nil {
//				log.Error("close关闭失败: %s", err.Error())
//			}
//			time.Sleep(time.Second * 5)
//			defer delete(manager.nodes, node.ClientId)
//			log.Info("心跳超时")
//			return
//		}
//	}
//}
