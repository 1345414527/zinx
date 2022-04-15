package ziface

/**
连接管理模块抽象层
*/
type IConnManager interface {
	//添加链接
	AddConn(conn IConnection)

	//删除链接
	RemoveConn(conn IConnection)

	//根据connID获取链接
	GetConn(connId uint32) (IConnection, error)

	//得到当前链接数目
	LenConn() int

	//清楚并终止所有的链接
	ClearConn()
}
