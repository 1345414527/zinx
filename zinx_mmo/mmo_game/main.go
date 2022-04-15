package main

import (
	"fmt"
	"zinxMmo/mmo_game/apis"
	"zinxMmo/mmo_game/core"
	"zinxMmo/zinx/ziface"
	"zinxMmo/zinx/znet"
)

/**
当前客户端建立连接之后的hook函数
*/
func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个Player对象
	player := core.NewPlayer(conn)

	//触发玩家上线
	player.Online()

	//将该链接绑定一个Pid
	conn.SetProperty("pid", player.Pid)

	fmt.Println("=====>Player pid = ", player.Pid, "is arrived <=====")
}

/**
当前客户端结束连接之前的hook函数
*/
func OnConnectionLost(conn ziface.IConnection) {
	//通过链接获取player
	pid, err := conn.GetProperty("pid")
	if err != nil {
		fmt.Println("Get LostPlayer Pid error", err)
		return
	}
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//触发玩家下线
	player.Offline()
}

func main() {
	//创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	//链接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	//启动服务
	s.Serve()

}
