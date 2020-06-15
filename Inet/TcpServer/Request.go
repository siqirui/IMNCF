package TcpServer

import (
	IinterFace "IMCF/IinterFace"
	_ "fmt"
)

type Request struct {
	conn IinterFace.IConnect //已经和客户端建立好的 链接
	data IinterFace.IMessage //客户端请求的数据
}

//获取请求连接信息
func (this *Request) GetConnection() IinterFace.IConnect {
	return this.conn
}

//获取请求消息的数据
func (this *Request) GetData() IinterFace.IMessage {
	return this.data
}

//获取事件类型
func (this *Request) GetMsgEvent() uint32 {
	return this.data.GetMsgEvent()
}
