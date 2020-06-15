package IinterFace

type IConnManager interface {
	AddConn(conn IConnect)                   //添加链接
	RemoveConn(conn IConnect)                //删除链接
	GetConn(connID uint32) (IConnect, error) //	根据链接ID获取链接
	GetConnNum() int                         //获取当前链接数量
	ClearConn()                              //清理所有链接
}
