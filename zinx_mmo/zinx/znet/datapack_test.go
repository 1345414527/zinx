package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟服务端
	*/
	//1.创建socketTCP
	listen, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//2.从客户端读取数据，拆包处理
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("server accept error:", err)
				continue
			}
			go func(conn net.Conn) {
				//处理客户端的请求
				dp := NewDataPack()
				for {
					//1. 第一次从head读，把包从head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData) //只获取head
					if err != nil {
						fmt.Println("read head error")
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						break
					}

					//2. 第二次从conn读，根据head中的datalen再读取data内容
					if msgHead.GetDataLen() > 0 {
						//msg是有效的，需要进行第二次读取
						//第二次从conn读，根据head中的datalen再读取data内容
						msg := msgHead.(*Message) //接口强转
						msg.SetData(make([]byte, msgHead.GetDataLen()))

						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.GetData())
						if err != nil {
							fmt.Println("server unpack data err", err)
							break
						}

						//完整的一个消息已经读取完毕
						fmt.Println("-->Recv MsgID:", msg.GetId(), "datalen = ", msg.DataLen, "data = ", string(msg.Data))

					}
				}

			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	//
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err", err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装连个msg一同发送
	//封装第一个msg包
	msg01 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("hello"),
	}

	sendData1, err := dp.Pack(msg01)
	if err != nil {
		fmt.Println("client pack msg01 err", err)
		return
	}

	//封装第二个msg包
	msg02 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte("nihaoya"),
	}

	sendData2, err := dp.Pack(msg02)
	if err != nil {
		fmt.Println("client pack msg02 err", err)
		return
	}

	//将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	//发送给服务器
	conn.Write(sendData1)

	//客户端阻塞
	select {}

}
