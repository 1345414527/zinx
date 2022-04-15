package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("client start")
	time.Sleep(time.Second)

	//1.直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}

	for {
		//发送封包的message消息
		dp := znet.NewDataPack()

		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx client Test Message")))
		if err != nil {
			fmt.Println("Pack err", err)
			return
		}
		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		//服务器就应该给我们回复一个message数据
		binaryDataHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryDataHead); err != nil {
			fmt.Println("read msg head error", err)
			return
		}

		fmt.Println(binaryDataHead)

		//先读取数据的head部分，得到id和datalen
		msg, err := dp.Unpack(binaryDataHead)
		if err != nil {
			fmt.Println("unpack err", err)
			return
		}

		//然后根据datalen进行第二次读取，得到data
		data := make([]byte, msg.GetDataLen())
		if msg.GetDataLen() > 0 {
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
		}
		msg.SetData(data)
		fmt.Println("--> recv server msg :ID = ", msg.GetId(), "len = ", msg.GetDataLen(), "data = ", string(msg.GetData()))

		time.Sleep(time.Second)
	}
}
