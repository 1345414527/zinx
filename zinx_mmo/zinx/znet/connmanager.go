package znet

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"zinxMmo/zinx/ziface"
)

/**
链接管理模块
*/

type ConnManager struct {
	//管理的链接集合
	connections map[uint32]ziface.IConnection

	//保护链接集合的读写锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//添加链接
func (cm *ConnManager) AddConn(conn ziface.IConnection) {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//将conn加入到集合中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connId", conn.GetConnID(), " add to ConnManager succ：conn num = ", cm.LenConn())
}

//删除链接
func (cm *ConnManager) RemoveConn(conn ziface.IConnection) {
	//加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//移除
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connId", conn.GetConnID(), "remove from ConnManager succ：conn num = ", cm.LenConn())
}

//根据connID获取链接
func (cm *ConnManager) GetConn(connId uint32) (ziface.IConnection, error) {
	//加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	}

	return nil, errors.New("不存在" + strconv.Itoa(int(connId)) + "对应的connection")
}

//得到当前链接数目
func (cm *ConnManager) LenConn() int {
	return len(cm.connections)
}

//清楚并终止所有的链接
func (cm *ConnManager) ClearConn() {
	//加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除conn并停止conn的工作
	for connId, conn := range cm.connections {
		//停止
		conn.Stop()

		//删除
		delete(cm.connections, connId)
	}

	fmt.Println("Clear All Connections secc! conn num = ", cm.LenConn())
}
