package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinxMmo/zinx/utils"
	"zinxMmo/zinx/ziface"
)

/**
封包，拆包的具体模块
使用的是TLV格式，物联网可以参考使用MQTT格式
*/

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的头的长度方法
func (d *DataPack) GetHeadLen() uint32 {
	return 8
}

//封包方法
//|datalen|msgID|data|
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//将MsgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetId()); err != nil {
		return nil, err
	}

	//将data数据写进dataBuff中

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法(将包的Head信息读出来)之后再根据head信息里的data长度，再进行一次读
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据二点ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息，得到daatlen和MSGID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && utils.GlobalObject.MaxPackageSize > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}
	return msg, nil
}
