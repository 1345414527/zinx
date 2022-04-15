package znet

type Message struct {
	Id      uint32 //消息的ID
	DataLen uint32 //消息的长度
	Data    []byte //消息的内容
}

//创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

//获取消息的id
func (m *Message) GetId() uint32 {
	return m.Id
}

//获取消息的长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

//获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

//设置消息的id
func (m *Message) SetId(id uint32) {
	m.Id = id
}

//设置消息的长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

//设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
