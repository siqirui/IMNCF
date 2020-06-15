package Utils

import (
	"IMCF/IinterFace"
	"config"
)

const (
	DATA_UPDATE_HEAD_SIZE = 4 + 1*8
)
const (
	DATA_UPDATE_BEFORE = 60
	DATA_UPDATE        = 61
	DATA_UPDATE_AFTER  = 62
)

type Global struct {
	TcpServer IinterFace.IServer `toml:"tcp_server"` //全局TCPServer
	Host      string             `toml:"host"`       //当前主机ip
	TcpPort   int                `toml:"tcp_port"`   //当前服务监听端口
	Name      string             `toml:"name"`       //当前服务名称

	Version          string `toml:"version"`             //当前服务版本
	MaxPackSize      uint32 `toml:"max_pack_size"`       //当前协议栈最大包
	MaxConn          int    `toml:"max_conn"`            //当前允许最大连接数
	WorkerPoolSize   uint32 `toml:"worker_pool_size"`    //业务任务池数量
	MaxWorkerTaskLen uint32 `toml:"max_worker_task_len"` //业务任务池 对应任务队列最大任务存储数量
	MaxMsgChanLen    uint32 `toml:"max_msg_chan_len"`    //sendbufmsg 发送消息的缓冲区大小
}

//定义全局配置
var GlobalOBJ *Global

//是否存在
//func PathExist(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//	return false, err
//}

func (this *Global) Reload(conf *config.Config) {

	GlobalOBJ.Host = conf.TcpServer.Host
	GlobalOBJ.TcpPort = conf.TcpServer.TcpPort
	GlobalOBJ.Name = conf.TcpServer.Name
	GlobalOBJ.Version = conf.TcpServer.Version
	GlobalOBJ.MaxPackSize = conf.TcpServer.MaxPackSize
	GlobalOBJ.MaxConn = conf.TcpServer.MaxConn
	GlobalOBJ.WorkerPoolSize = conf.TcpServer.WorkerPoolSize
	GlobalOBJ.MaxWorkerTaskLen = conf.TcpServer.MaxWorkerTaskLen
	GlobalOBJ.MaxMsgChanLen = conf.TcpServer.MaxMsgChanLen

}

//默认
func Init() {
	GlobalOBJ = &Global{
		Host:             "0.0.0.0",
		TcpPort:          3030,
		Name:             "IMFCServer",
		Version:          "V0.0.1",
		MaxPackSize:      10240,
		MaxConn:          120400,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 512,
		MaxMsgChanLen:    120,
	}
	//优先配置文件中的配置信息
	GlobalOBJ.Reload(config.GetConf())

}
