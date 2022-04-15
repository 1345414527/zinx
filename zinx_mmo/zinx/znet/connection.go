package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinxMmo/zinx/utils"
	"zinxMmo/zinx/ziface"
)

type Connection struct {
	//当前链接所属的server
	server ziface.IServer

	//当前链接两点socket tcp套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的链接状态
	isClosed bool

	//告知当前链接已经退出的/停止 channel
	ExitChan chan bool

	//无缓冲的管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	//该链接处理的方法Router
	MsgHandler ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}

	//保护链接属性的锁
	propertyLock sync.RWMutex
}

//初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		server:     server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
		property:   make(map[string]interface{}),
	}

	//将当前链接保存到链接管理器中
	server.GetConnManager().AddConn(c)

	return c
}

//链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit,remote addr is", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大512B
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 8B
		headData := make([]byte, dp.GetHeadLen())

		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包，得到msgID 和 msgDatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据datalen再次读取Data，放在msg.Data中
		data := make([]byte, msg.GetDataLen())
		if msg.GetDataLen() > 0 {
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		//判断是否开启工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//交给线程池
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			//执行注册的路由方法
			c.MsgHandler.DoMsgHandler(req)
		}

		//从路由中，找到注册绑定的Conn对应的router调用

	}

}

//写消息，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Write  Goroutine is running...]")
	defer fmt.Println("connID=", c.ConnID, "Writer is exit,remote addr is", c.GetRemoteAddr().String())

	//不断地阻塞等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error,", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要推出
			return
		}

	}

}

//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	//启动从当前链接的读数据的业务
	go c.StartReader()
	//启动从当前链接写数据的业务
	go c.StartWriter()

	//此时已经创建链接成功，调用用户创建的Hook函数
	c.server.CallOnConnStart(c)

}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	//当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//移除链接
	c.server.GetConnManager().RemoveConn(c)

	//此时链接将要关闭，调用用户创建的Hook函数
	c.server.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	//告知writer关闭
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

//获取当前链接的绑定socket conn
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供一个sendMSG方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))

	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个属性
	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	//获取属性
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New(" no property found!")
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//删除属性
	delete(c.property, key)
}
