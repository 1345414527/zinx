package core

import (
	"sync"
)

/**
当前游戏的世界总管理模块
*/

type WorldManager struct {
	//AOIManager 当前世界地图AOI的管理模块
	AoiMgr *AOIManager

	//当前全部在线的Players集合
	Players map[int32]*Player

	//保护Players集合的锁
	pLock sync.RWMutex
}

//提供一个全局的世界管理模块的句柄
var WorldMgrObj *WorldManager

/**
初始化方法
*/
func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_X, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}

/**
添加一个玩家
*/
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	//将player添加到AOIManager中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

/**
删除一个玩家
*/
func (wm *WorldManager) DeleterPlayer(pid int32) {
	//通过pid获取玩家
	player := wm.Players[pid]
	wm.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Z)

	wm.pLock.Lock()
	delete(wm.Players, pid)
	wm.pLock.Unlock()

}

/**
通过玩家ID查询Player玩家
*/
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	player := wm.Players[pid]
	wm.pLock.RUnlock()
	return player
}

/**
获取全部的在线玩家
*/
func (wm *WorldManager) GetAllOnlinePlayer() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)
	for _, p := range wm.Players {
		players = append(players, p)
	}

	return players
}
