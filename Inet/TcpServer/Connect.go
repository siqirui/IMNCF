package TcpServer

import (
	IinterFace "IMCF/IinterFace"
	"IMCF/Utils"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connect struct {
	TcpServer IinterFace.IServer //当前conn归属的server
	Conn      *net.TCPConn       //当前连接的socket TCP套接字
	ConnID    uint32             //当前连接的ID 也可以称作为SessionID，ID全局唯一
	isClosed  bool               //当前连接的关闭状态
	//handleFunc   IinterFace.HandFunc //该连接的处理方法api
	ExitBuffChan chan bool //告知该链接已经退出/停止的channel
	//Route        IinterFace.IRoute
	msgHandles IinterFace.IMsgHandler //msgEvent 对应处理方法
	msgChan    chan []byte            //无缓冲管道，用于读写两个协成通讯
	msgBufChan chan []byte            //有缓冲消息管道

	property map[string]interface{} //链接属性

	propertyLock sync.RWMutex //保护链接属性读写锁

	//给缓冲队列发送数据的channel，
	// 如果向缓冲队列发送数据，那么把数据发送到这个channel下
	//	SendBuffChan chan []byte
}

func NewConnect(server IinterFace.IServer, conn *net.TCPConn, connectID uint32, handles IinterFace.IMsgHandler) *Connect {
	connect := &Connect{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connectID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		msgHandles:   handles,
		msgChan:      make(chan []byte),
		msgBufChan:   make(chan []byte, Utils.GlobalOBJ.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}
	connect.TcpServer.GetConnManager().AddConn(connect)

	return connect

}

func (this *Connect) StartWriter() {
	defer fmt.Println(this.RemoteAddr().String(), "conn writer exit!")
	defer this.Stop()
	for {
		select {
		case data := <-this.msgChan:
			if _, err := this.Conn.Write(data); err != nil {
				fmt.Println("Send Data error :,", err, "conn writer exit")
				return
			}
		case data, ok := <-this.msgBufChan:
			if ok {
				if _, err := this.Conn.Write(data); err != nil {
					fmt.Println("Send Buf Data error :,", err, "conn writer exit")
					return
				}
			} else {
				break
				fmt.Println("[ERROR] msgBufChan is closed")
			}
		case <-this.ExitBuffChan:
			return
		}
	}
}

func (this *Connect) StartRead() {
	fmt.Println("Read Goroutine running")
	defer fmt.Println(this.RemoteAddr().String(), " conn reader exit!")
	defer this.Stop()

	for {
		//buf := make([]byte, 1024)
		//_, err := this.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err : ", err)
		//	this.ExitBuffChan <- true
		//	continue
		//}
		//dataPack DataPack IDataPack

		dataPack := NewDataPack()
		headData := make([]byte, dataPack.GetHeadLen())
		if _, err := io.ReadFull(this.GetTcpConnect(), headData); err != nil {
			fmt.Println("[ERROR]read msg head error ", err)
			break
		}

		//msg IMessage
		msg, err := dataPack.UnPack(headData)
		if err != nil {
			fmt.Println("[ERROR]unpack error ", err)
			break
		}
		if msg.GetMsgEvent() == Utils.DATA_UPDATE_AFTER {
			fmt.Println("recv buf -->60 : ", headData)
			msg.SetDataLen(0)
		} else {
			datalenbyte := make([]byte, 4)
			if _, err := io.ReadFull(this.GetTcpConnect(), datalenbyte); err != nil {
				fmt.Println("[ERROR]read msg len error ", err)
				break
			} else {
				dataLen := binary.BigEndian.Uint32(datalenbyte)
				msg.SetDataLen(dataLen)
				fmt.Println("recv data len ", dataLen)
			}
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(this.GetTcpConnect(), data); err != nil {
				fmt.Println("[ERROR]read msg data error ", err)

				this.ExitBuffChan <- true

				break
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		var req = Request{
			conn: this,
			data: msg,
		}
		//从路由Routers 中找到注册绑定Conn的对应Handle

		if Utils.GlobalOBJ.WorkerPoolSize > 0 {
			//如果工作池已启动 则交给工作池去处理
			this.msgHandles.SendMsgToTaskQueue(&req)
		} else {
			//如果工作池没有启动，则交给绑定好的handle处理
			go this.msgHandles.DoMsgHandler(&req)
		}
		//go this.msgHandles.DoMsgHandler(&req)
		//go func(request IinterFace.IRequest) {
		//	//执行注册的路由方法
		//	this.Route.PreHandle(request)
		//	this.Route.Handle(request)
		//	this.Route.PostHandle(request)
		//}(&req)
		//HandFunc 	func(*net.TCPConn, []byte, int) error
		//if err := this.handleFunc( this.Conn , buf, len);err != nil{
		//if err := this.handleFunc(this.Conn, buf, len); err !=nil {
		//	fmt.Println("do handleFunc err : ",err,"connectID : ",this.ConnID)
		//	this.ExitBuffChan <- true
		//}
	}
}

func (this *Connect) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	go this.StartRead()
	go this.StartWriter()

	this.TcpServer.CallOnConnStart(this)
	//for {
	//	select {
	//	case <-this.ExitBuffChan:
	//		return
	//	}
	//}
}

func (this *Connect) Stop() {
	//1. 如果当前链接已经关闭
	if this.isClosed == true {
		return
	}
	this.isClosed = true
	this.TcpServer.CallOnConnStop(this)
	this.Conn.Close()

	this.ExitBuffChan <- true

	this.TcpServer.GetConnManager().RemoveConn(this)
	//关闭该链接全部管道
	close(this.ExitBuffChan)
	close(this.msgBufChan)
}

func (this *Connect) GetTcpConnect() *net.TCPConn {
	return this.Conn
}

func (this *Connect) GetTcpConnectID() uint32 {
	return this.ConnID
}

func (this *Connect) RemoteAddr() net.Addr {
	return this.Conn.RemoteAddr()
}

func (this *Connect) SendMsg(msgEvent uint32, fileID [8]byte, data []byte) error {
	if this.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dataPack := NewDataPack()
	msg, err := dataPack.Pack(NewMsgPackage(msgEvent, fileID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgEvent)
		return errors.New("Pack error msg ")
	}

	fmt.Println("send buf :", msg)
	//写回客户端
	this.msgChan <- msg
	return nil
}

//func (this *Connect) Send(data []byte) error {
//	return nil
//}
//

func (this *Connect) SendBufMsg(msgEvent uint32, fileID [8]byte, data []byte) error {
	if this.isClosed {
		return errors.New("Connect closed when send buf msg")
	}
	dataPack := NewDataPack()
	msg, err := dataPack.Pack(NewMsgPackage(msgEvent, fileID, data))
	if err != nil {
		fmt.Println("[ERROR] Pack error msg event = ", msgEvent)
		return errors.New("Pack error msg")
	}
	this.msgBufChan <- msg
	return nil
}
func (this *Connect) SetProperty(key string, value interface{}) {
	this.propertyLock.Lock()
	defer this.propertyLock.Unlock()

	this.property[key] = value
}
func (this *Connect) GetProperty(key string) (interface{}, error) {
	this.propertyLock.Lock()
	defer this.propertyLock.Unlock()
	if value, ok := this.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("property no found")
	}
}
func (this *Connect) RemoveProperty(key string) {
	this.propertyLock.Lock()
	defer this.propertyLock.Unlock()

	delete(this.property, key)
}
