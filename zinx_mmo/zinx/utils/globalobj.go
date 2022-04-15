package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinxMmo/zinx/ziface"
)

/*
存储一切有关Zinx框架的全局参数，供其它模块使用
一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
		Zinx
	*/
	Version          string //当前Zinx的版本号
	MaxConn          int    //当前服务器主机允许的最大链接数
	MaxPackageSize   uint32 //当前zinx框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //Zinx框架运行用户最多开辟多少个Worker(限定条件)
}

/*
定义一个全局的对外Globalobj
*/
var GlobalObject *GlobalObj

/*
从zinx.json中加载用于自定义的参数
*/
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("demoTest/ZinxV0.2/conf/zinx.json")
	if err != nil {
		fmt.Println("reload zinx.json err", err)
		panic(err)
	}
	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		fmt.Println("resolve zinx.json to GlobalObject err", err)
		panic(err)
	}
}

/*
提供一个init方法，初始化当前的GlobalObj
*/
func init() {
	//如果配置文件没有加载就是默认的一个值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.4",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   1024 * 4,
		WorkerPoolSize:   10,   //工作池中队列的个数
		MaxWorkerTaskLen: 1024, //每个池子对应的消息队列的任务的数量最大值
	}

	//从conf/zinx.json中加载用户自定义的参数
	//GlobalObject.Reload()
}
