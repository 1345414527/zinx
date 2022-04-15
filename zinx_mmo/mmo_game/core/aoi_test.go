package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(100, 300, 200, 450, 4, 5)
	fmt.Println(aoiMgr)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 0, 250, 5, 5)

	for gid, _ := range aoiMgr.grids {
		//得到当前gid的周边九宫格信息
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid:", gid, "grids len = ", len(grids))
		gIds := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIds = append(gIds, grid.Gid)
		}
		fmt.Println("surrounding grid IDs are ", gIds)
	}
}
