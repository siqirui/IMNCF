package IMCF

/**

*  IMCFV0.1，测试
 */

import (
	"IMCF/IinterFace"
	_ "IMCF/Ilog"
	"IMCF/Inet/TcpServer"
	utils "IMCF/Utils"
	"IMCF/Worker"
	"logging"
)

//创建连接的时候执行
func DoConnectionBegin(conn IinterFace.IConnect) {
	logging.Logger.Debug("DoConnecionBegin is Called ... ")
	logging.Logger.Debug("new conn ip : %s", conn.GetTcpConnect().LocalAddr())
}

//连接断开的时候执行
func DoConnectionLost(conn IinterFace.IConnect) {
	//在连接销毁之前，查询conn的Name，Home属性
	//delete(utils.ConnectFiles,conn.GetTcpConnectID())
	if name, err := conn.GetProperty("Name"); err == nil {
		logging.Logger.Error("Conn Property Name = ", name)
		//log.Error("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		//Ilog.Error("Conn Property Home = ", home)
		logging.Logger.Error("Conn Property Home = ", home)
	}
	conn.Stop()
	//Ilog.Debug("DoConneciotnLost is Called ... ")
	logging.Logger.Debug("DoConneciotnLost is Called ... ")
}

func TcpServerStart() {
	//创建一个server句柄
	server := TcpServer.NewServer()

	//注册链接hook回调函数
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)

	//配置路由

	server.AddRouter(utils.DATA_UPDATE_BEFORE, &Worker.Router{})
	server.AddRouter(utils.DATA_UPDATE, &Worker.Router{})
	server.AddRouter(utils.DATA_UPDATE_AFTER, &Worker.Router{})

	//开启服务
	server.Server()
}
