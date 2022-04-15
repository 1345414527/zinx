package znet

import (
	"fmt"
	"strconv"
	"zinxMmo/zinx/utils"
	"zinxMmo/zinx/ziface"
)

/**
消息处理模块的实现
*/
type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter

	//负责Worker取任务的消息队列
	TaskQueue chan ziface.IRequest

	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

//初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      nil,
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, flag := mh.Apis[msgID]; flag {
		//已经存在
		panic("repeat api,msgID = " + strconv.Itoa(int(msgID)))
	}
	//添加
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, "succ!")
}

/**
单个进行处理
*/

//调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//从request中找到msgID
	router, flag := mh.Apis[request.GetMsg().GetId()]
	if !flag {
		fmt.Println("api MsgID = ", request.GetMsg().GetId(), " is NOT FOUND!Need Register!")
		return
	}
	//根据MsgID调度对应的router业务即可
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

/**
使用线程池来进行处理任务
*/

//启动一个worker工作池(只发生一次,一个zinx只有一个工作池，在server中创建)
func (mh *MsgHandle) StartWorkerPool() {
	//创建一个消息队列
	mh.TaskQueue = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
	//创建多个线程进行监听
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		id := i
		fmt.Println("Worker ID = ", id, "is started...")
		go func() {
			for {
				select {
				case request := <-mh.TaskQueue:
					fmt.Println(id, " 正在处理：")
					mh.DoMsgHandler(request)
				}
			}
		}()
	}
}

//发送给channel
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//发送给TaskQueue
	mh.TaskQueue <- request
}
