package ziface

/**
将请求的消息封装到一个Message中，定义抽象接口
*/
type IMessage interface {
	//获取消息的id
	GetId() uint32

	//获取消息的长度
	GetDataLen() uint32

	//获取消息的内容
	GetData() []byte

	//设置消息的id
	SetId(id uint32)

	//设置消息的长度
	SetDataLen(len uint32)

	//设置消息的内容
	SetData(data []byte)
}
