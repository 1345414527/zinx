package core

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
	"zinxMmo/mmo_game/pb"
	"zinxMmo/zinx/ziface"
)

/**
玩家对象
*/

type Player struct {
	Pid  int32              //玩家id
	Conn ziface.IConnection //当前玩家的连接
	X    float32            //平面x坐标
	Y    float32            //高度
	Z    float32            //平面y坐标
	V    float32            //旋转的0-360°
}

/*
Player ID 生成器（一般都是从数据库获取）
*/
var (
	PidGen int32 = 1
	IdLock sync.RWMutex
)

/**
创建一个玩家的方法
*/
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(134 + rand.Intn(10)),
		V:    0,
	}

	fmt.Printf("New a player pid: %d,X: %v,Y %v,Z: %v,V %v\n", p.Pid, p.X, p.Y, p.Z, p.V)

	return p
}

/**
提供一个发送给客户端消息的方法
主要是将pb的protobuf数据序列化之后，再调用zinx的SendMsg方法
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将proto Message结构体序列化 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err", err)
		return
	}

	//将二进制文件 通过zinx框架的sendmsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil", err)
		return
	}

	if err = p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error!", err)
		return
	}
}

/**
玩家上线
*/

func (p *Player) Online() {
	//给客户端发送MsgID：1的消息:同步当前player的id给客户端
	p.SyncPid()

	//给客户端发送MsgID：200的消息：同步当前player的初始位置给客户端
	p.BroadCastStartPosition()

	//将当前新上线的玩家添加到WorldManager中
	WorldMgrObj.AddPlayer(p)

	//同步周边玩家，告知他们当前玩家已经上线，广播当前玩家的位置信息
	p.SyncSurrounding()
}

/**
玩家下线
*/
func (p *Player) Offline() {
	//得到当前玩家周边的九宫盒内有哪些玩家
	players := p.GetSurroundingPlayers()

	//给周围玩家广播MSGID:201消息
	proto_msg := &pb.SyncPID{
		PID: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.DeleterPlayer(p.Pid)

}

/**
告知玩家pid，同步已经已经生成的玩家ID给客户端
*/
func (p *Player) SyncPid() {
	//组建MsgID：0 的proto数据
	proto_msg := &pb.SyncPID{
		PID: p.Pid,
	}

	//将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

/**
广播玩家自己的出生地点
*/
func (p *Player) BroadCastStartPosition() {
	//组件MsgID: 200的proto数据
	proto_msg := &pb.BroadCast{
		PID: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
			},
		},
	}

	//将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

/**
玩家广播世界聊天信息
*/
func (p *Player) Talk(content string) {
	//1. 组建MsgID：200 proto数据
	proto_msg := &pb.BroadCast{
		PID: p.Pid,
		Tp:  1, //1代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//2. 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllOnlinePlayer()

	//3. 向所有的玩家发送消息MsgID: 200
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

/**
获取当前玩家的九宫格玩家
*/
func (p *Player) GetSurroundingPlayers() []*Player {
	//获取九宫格里的pid
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	//根据pid获取player
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}

/**
同步周边玩家
*/
func (p *Player) SyncSurrounding() {
	//1. 获取当前玩家的九宫格
	players := p.GetSurroundingPlayers()

	//2. 将当前玩家的位置信息通过MsgID：200发送给周围玩家（让周围玩家看到自己）
	//组件MsgID：200 proto数据
	proto_msg := &pb.BroadCast{
		PID: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//分别向周围的玩家发送消息
	for _, player := range players {
		fmt.Println("发送消息", player.Pid)
		player.SendMsg(200, proto_msg)
	}

	//3. 将周围的全部玩家的位置信息发送给给当前的玩家MsgID：202 客户端（让自己看到其它玩家）
	//组建MsgID：202 proto数据
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			PID: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}

	syncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg,
	}

	//将组建好的数据发送给当前玩家的客户端
	p.SendMsg(202, syncPlayers_proto_msg)
}

/**
广播当前玩家的位置移动信息
*/
func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//组件广播proto协议 MsgID：200 TP-4
	proto_msg := &pb.BroadCast{
		PID: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: x,
				Y: y,
				Z: z,
				V: v,
			},
		},
	}
	//获取当前玩家的周边玩家
	players := p.GetSurroundingPlayers()
	//一次给所有玩家对应的客户端发送更新信息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

}
