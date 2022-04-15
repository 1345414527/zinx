package core

import (
	"fmt"
	"sync"
)

/**
一个AOI地图中的格子
*/

type Grid struct {
	//格子ID
	Gid int

	//格子的左边边界坐标
	MinX int

	//格子的右边边界坐标
	MaxX int

	//格子的上边边界坐标
	MinY int

	//格子的下边边界坐标
	MaxY int

	//当前格子内玩家或者物体成员的ID集合
	playerIds map[int]bool

	//保护当前集合的锁
	pIdLock sync.RWMutex
}

//初始化方法
func NewGrid(gId, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		Gid:       gId,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIds: make(map[int]bool),
	}
}

//给格子添加一个玩家
func (g *Grid) Add(playerId int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()
	g.playerIds[playerId] = true

}

//从格子中删除一个玩家
func (g *Grid) Remove(playerId int) {
	g.pIdLock.Lock()
	defer g.pIdLock.Unlock()

	delete(g.playerIds, playerId)
}

//得到当前格子中所有的玩家
func (g *Grid) GetPlayerIds() (playerIds []int) {
	g.pIdLock.RLock()
	defer g.pIdLock.RUnlock()

	for k, _ := range g.playerIds {
		playerIds = append(playerIds, k)
	}

	return
}

//调试使用-打印出格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id : %d,minX: %d,maxX: %d,minY: %d,maxY: %d,playerIds: %v\n",
		g.Gid, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIds)
}
