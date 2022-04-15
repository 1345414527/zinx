package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinxMmo/mmo_game/core"
	"zinxMmo/mmo_game/pb"
	"zinxMmo/zinx/ziface"
	"zinxMmo/zinx/znet"
)

/**
世界聊天 路由业务
*/

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	//1. 解析客户端传递过来的proto协议
	proto_msg := new(pb.Talk)
	if err := proto.Unmarshal(request.GetMsg().GetData(), proto_msg); err != nil {
		fmt.Println("Talk Unmarshal error", err)
		return
	}

	//2. 当前的聊天数据是属于哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("Get Talker Pid error", err)
		return
	}

	//3. 根据pid的得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4. 将这个消息广播给其它全部在线的玩家
	player.Talk(proto_msg.Content)
}
