package znet

import (
	"fmt"
	"net"
	"zinxMmo/zinx/utils"
	"zinxMmo/zinx/ziface"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string

	//服务器绑定的ip地址
	IPVersion string

	//服务器坚挺的ip
	IP string

	//服务器监听的端口
	Port int

	//当前server的消息管理模块，用来绑定MsgID和对应的处理业务的API关系
	MsgHandler ziface.IMsgHandle

	//该server的链接管理器
	connManager ziface.IConnManager

	//该Server创建链接之后自动调用Hook函数-OnConnStart
	OnConnStart func(conn ziface.IConnection)

	//该Server销毁链接之前自动调用Hook函数-OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

/**
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		connManager: NewConnManager(),
	}
	s.Name = name

	return s
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s,Port %d, is starting\n", s.IP, s.Port)

	go func() {
		//开启线程池
		s.MsgHandler.StartWorkerPool()

		//1.解析一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error：", err)
			return
		}
		//2.获取监听器
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}

		fmt.Println("start Zinx server succ", s.Name, " succ ,Listenning: ")
		var cid uint32
		cid = 0

		//3.阻塞的等待客户端链接，处理客户端链接业务
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大连接数的判断
			if s.connManager.LenConn() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大链接的错误包
				conn.Close()
				continue
			}

			//将处理新连接的业务方法 和 conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	//将一些服务器的资源、状态或者一些已经开辟的链接信息进行停止或回收
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.connManager.ClearConn()

}

//运行服务器
func (s *Server) Serve() {
	//启动server
	s.Start()

	//做一些服务器启动后的额外业务

	//阻塞状态
	select {}
}

//添加router
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Succ!!")
}

//获取链接管理器
func (s *Server) GetConnManager() ziface.IConnManager {
	return s.connManager
}

//注册OnConnStart钩子函数
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子函数
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子函数
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart()==>")
		s.OnConnStart(connection)
	}
}

//调用OnConnStop钩子函数
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop()==>")
		s.OnConnStop(connection)
	}
}
