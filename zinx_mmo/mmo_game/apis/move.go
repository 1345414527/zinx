package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinxMmo/mmo_game/core"
	"zinxMmo/mmo_game/pb"
	"zinxMmo/zinx/ziface"
	znet2 "zinxMmo/zinx/znet"
)

/**
玩家移动路由
*/

type MoveApi struct {
	znet2.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	//解析客户端传递过来的proto协议
	proto_msg := &pb.Position{}
	if err := proto.Unmarshal(request.GetMsg().GetData(), proto_msg); err != nil {
		fmt.Println("Move: Position Unmarshal error", err)
		return
	}

	//得到当前发送位置的是哪个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		return
	}
	fmt.Printf("Player pid = %d,move (%f,%f,%f,%f)", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	//给其他玩家进行当前玩家的位置信息广播
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//广播并更新
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
