package IinterFace

type IRequest interface {
	GetConnection() IConnect //获取请求连接信息
	GetData() IMessage       //获取请求消息的数据
	GetMsgEvent() uint32     //获取请求的消息ID
}
