package IinterFace

//服务器接口

type IServer interface {
	Start()                                    //启动服务
	Stop()                                     //停止服务
	Server()                                   //启动业务
	AddRouter(msgEvent uint32, router IRouter) //路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	GetConnManager() IConnManager              //获取一个链接管理器
	SetOnConnStart(func(connect IConnect))     //设置该Server的链接创建时hook
	SetOnConnStop(func(connect IConnect))      //设置该Server的链接断开时hook
	CallOnConnStart(conn IConnect)             //调用OnConnStart hook
	CallOnConnStop(conn IConnect)              //调用OnConnStop hook
}
