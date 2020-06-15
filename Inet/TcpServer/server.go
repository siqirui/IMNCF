package TcpServer

import (
	IinterFace "IMCF/IinterFace"
	Utils "IMCF/Utils"
	"config"
	"fmt"
	"net"
	"time"
)

//实现一个server类
type Server struct {
	Name        string                         //服务名称
	IpVersion   string                         //IP版本
	Ip          string                         //服务器地址
	Port        int                            //服务器端口
	msgHandler  IinterFace.IMsgHandler         //消息管理模块
	ConnManager IinterFace.IConnManager        //链接管理模块
	OnConnStart func(conn IinterFace.IConnect) //该server的链接创建时hook
	OnConnStop  func(conn IinterFace.IConnect) //断开时的hook
}

//创建server实例

func NewServer() IinterFace.IServer {
	Utils.GlobalOBJ.Reload(config.GetConf())
	server := &Server{
		Name:        Utils.GlobalOBJ.Name,
		IpVersion:   "tcp4",
		Ip:          Utils.GlobalOBJ.Host,
		Port:        Utils.GlobalOBJ.TcpPort,
		msgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return server
}

//相关接口实现
func (this *Server) Start() {
	fmt.Printf("[START]Server Name : %s,listener at IP %s,Port:%d is starting \n",
		this.Name, this.Ip, this.Port)
	fmt.Printf("[IMCF] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		Utils.GlobalOBJ.Version,
		Utils.GlobalOBJ.MaxConn,
		Utils.GlobalOBJ.MaxPackSize)
	go func() {
		//0启动任务池
		this.msgHandler.StartWorkerPool()
		//1获取TCP
		addr, err := net.ResolveTCPAddr(this.IpVersion, fmt.Sprintf("%s:%d", this.Ip, this.Port))
		if err != nil {
			fmt.Println("do ResolveTCPAddr error: ", err)
			return
		}
		//2启动监听
		listener, err := net.ListenTCP(this.IpVersion, addr)
		if err != nil {
			fmt.Println("do ListenTCP err : ", err)
			return
		}
		//TODO需要一个生成cid的方法
		var cid uint32
		cid = 0
		//3启动网络
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("do AcceptTCP error : ", err)
				time.Sleep(1)
				continue
			}

			if this.ConnManager.GetConnNum() >= Utils.GlobalOBJ.MaxConn {
				conn.Close()
				continue
			}
			//fmt.Printf("conn succeed! remote addr : ", conn.RemoteAddr().String())
			//tmpData := make([]byte, 1024)
			//readlen, err := conn.Read(tmpData)
			//fmt.Printf("recv len :%d \n recv data : %s", readlen, tmpData)
			iconn := NewConnect(this, conn, cid, this.msgHandler)
			if iconn == nil {
				fmt.Println("[ERROR] NewConnect Error")
				return
			}
			cid++
			//启动当前连接处理业务
			go iconn.Start()

		}
	}()
}

func (this *Server) Stop() {
	fmt.Println("[STOP]IMFC server,name : ", this.Name)
	this.ConnManager.ClearConn()
}

func (this *Server) Server() {
	this.Start()
	select {}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (this *Server) AddRouter(msgEvent uint32, router IinterFace.IRouter) {
	this.msgHandler.AddRouter(msgEvent, router)
	fmt.Println("[INFO]Add Router succ! ")
}

func (this *Server) GetConnManager() IinterFace.IConnManager {
	return this.ConnManager
}
func (this *Server) SetOnConnStart(hookFunc func(connect IinterFace.IConnect)) {
	this.OnConnStart = hookFunc
}
func (this *Server) SetOnConnStop(hookFunc func(connect IinterFace.IConnect)) {
	this.OnConnStop = hookFunc
}
func (this *Server) CallOnConnStart(connect IinterFace.IConnect) {
	if this.OnConnStart != nil {
		fmt.Println("[INFO] wilL CallOnConnStart")
		this.OnConnStart(connect)
	}
}
func (this *Server) CallOnConnStop(connect IinterFace.IConnect) {
	if this.OnConnStop != nil {
		fmt.Println("[INFO] wilL CallOnConnStop")
		this.OnConnStop(connect)
	}

}
