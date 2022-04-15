package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
	基于zinx框架来开发的服务器端应用程序
*/

//继承BaseRouter
type PingRouter struct {
	znet.BaseRouter
}

//Test preHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("	Call Router PreHandle...")

}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("	Call Ping Router Handle...")
	//先读取客户端的数据,再写回ping...ping...ping
	fmt.Println("	recv from client: msgID = ", request.GetMsg().GetId(), "len = ", request.GetMsg().GetDataLen(), ",data = ", string(request.GetMsg().GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}

}

//Test postHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("	Call Router PostHandle...")

}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("	Call Hello Router Handle...")
	//先读取客户端的数据,再写回ping...ping...ping
	fmt.Println("	recv from client: msgID = ", request.GetMsg().GetId(), "len = ", request.GetMsg().GetDataLen(), ",data = ", string(request.GetMsg().GetData()))
	err := request.GetConnection().SendMsg(200, []byte("hello...hello...hello"))
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	//创建一个server
	ser := znet.NewServer("[zine v0.4]")

	//添加一个自定义router
	ser.AddRouter(0, &PingRouter{})
	ser.AddRouter(1, &HelloZinxRouter{})
	//启动server
	ser.Serve()
}
