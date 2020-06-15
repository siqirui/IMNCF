package IinterFace

import (
	_ "fmt"
	"net"
)

type IConnect interface {
	Start()                      //启动当前连接
	Stop()                       //关闭当前连接
	GetTcpConnect() *net.TCPConn //获取当前TCP连接句柄
	GetTcpConnectID() uint32     //获取当前TCP连接id
	RemoteAddr() net.Addr        //获取远程客户端地址
	//Send(data []byte) error      //无缓冲发送数据
	//SendBuf(data []byte) error   //有缓冲发送数据
	SendMsg(msgEvent uint32, fileID [8]byte, data []byte) error    //直接将Message数据发送数据给远程的TCP客户端
	SendBufMsg(msgEvent uint32, fileID [8]byte, data []byte) error //有缓冲
	SetProperty(key string, value interface{})                     //设置链接属性
	GetProperty(key string) (interface{}, error)                   //获取链接属性
	RemoveProperty(key string)                                     //删除链接属性
}
