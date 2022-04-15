package core

import (
	"fmt"
)

/*
定义一些AOI的边界值
*/
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_X int = 10
	AOI_CNTS_Y int = 20
)

/**
格子管理模块
*/

type AOIManager struct {
	//区域的左边界坐标
	MinX int

	//区域的右边界坐标
	MaxX int

	//区域的上边界坐标
	MinY int

	//区域的下边界坐标
	MaxY int

	//x方向的格子数量
	CntsX int

	//y方向的格子数量
	CntsY int

	//当前区域中有哪些格子map-key = 格子的ID，value = 格子对象
	grids map[int]*Grid
}

/**
初始化一个AOI区域的管理模块
*/
func NewAOIManager(minX, maxX, minY, maxY, cntsX, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		CntsX: cntsX,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	//给AOI初始化区域的格子所有的格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid := y*cntsX + x
			aoiMgr.grids[gid] = NewGrid(gid, aoiMgr.MinX+x*aoiMgr.gridWidth(), aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridHeight(), aoiMgr.MinY+(y+1)*aoiMgr.gridHeight())
		}
	}

	return aoiMgr
}

//得到每个格子在x轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个格子在y轴方向的长度
func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//根据格子GID得到周边九宫格格子的ID集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//判断gid是否在aoimanager中
	if _, ok := m.grids[gID]; !ok {
		return
	}

	//初始化grids切片
	grids = append(grids, m.grids[gID])

	//需要gId的左边是否有格子？右边是否有格子
	//需要通过gId得到当前格子在x轴的编号
	idx := gID % m.CntsX
	//判断idx编号是否有左格子，有则放入grids集合中
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	//判断idx编号是否有右格子，有则放入grids集合中
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	//获得gidsx
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.Gid)
	}

	//遍历gidsX集合中每个gid
	for _, v := range gidsX {
		idy := v / m.CntsY
		//gid上边是否有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])

		}
		//gid下边是否有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}
	return
}

//通过x，y坐标获取当前的grid编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridHeight()
	return idy*m.CntsX + idx
}

//通过x，y得到九宫格内全部的PlayerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIds []int) {
	//得到当前格子的id
	gId := m.GetGidByPos(x, y)

	//通过gId得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gId)
	fmt.Println("Nine grids info：", grids)

	//将九宫格的信息的全部的player的id累加到playerIds
	for _, grid := range grids {
		playerIds = append(playerIds, grid.GetPlayerIds()...)
		fmt.Printf("==> grid ID: %d,pids:%v===", grid.Gid, grid.GetPlayerIds())
	}

	return
}

//添加一个PlayerId到一个格子中
func (m *AOIManager) AddPidToGrid(pId, gId int) {
	m.grids[gId].Add(pId)
}

//移除一个格子中的PlayerId
func (m *AOIManager) RemovePidFromGrid(pId, gId int) {
	m.grids[gId].Remove(pId)
}

//通过gId获取全部的PlayerId
func (m *AOIManager) GetPidsByGid(gId int) (playerIds []int) {
	return m.grids[gId].GetPlayerIds()
}

//通过坐标将Player添加在一个格子中
func (m *AOIManager) AddToGridByPos(pId int, x, y float32) {
	m.AddPidToGrid(pId, m.GetGidByPos(x, y))
}

//通过坐标将Player从一个格子中删除
func (m *AOIManager) RemoveFromGridByPos(pId int, x, y float32) {
	m.RemovePidFromGrid(pId, m.GetGidByPos(x, y))
}

//打印格子的信息
func (m *AOIManager) String() string {

	s := fmt.Sprintf("AOIManager:\n\tMinx:%d,MaxX:%d,MinY:%d,MaxY:%d,cntsX:%d,cntsY:%d\n\tGrids in AOIManager：\n",
		m.MinX, m.MaxX, m.MinY, m.MaxY, m.CntsX, m.CntsY)

	//打印全部格子信息
	for _, grid := range m.grids {
		s += fmt.Sprint("\t\t")
		s += fmt.Sprintln(grid)
	}
	return s
}
