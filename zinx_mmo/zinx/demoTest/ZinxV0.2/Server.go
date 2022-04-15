package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
	基于zinx框架来开发的服务器端应用程序
*/

type PingRouter struct {
	znet.BaseRouter
}

//Test preHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("ping ping..."))
	if err != nil {
		fmt.Println("call back  ping ping error")
	}
}

//Test postHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTcpConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	//创建一个server
	ser := znet.NewServer("[zine v0.2]")

	//添加一个自定义router
	ser.AddRouter(0, &PingRouter{})
	//启动server
	ser.Serve()
}
