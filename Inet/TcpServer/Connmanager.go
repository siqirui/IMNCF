package TcpServer

import (
	IinterFace "IMCF/IinterFace"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connects map[uint32]IinterFace.IConnect
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connects: make(map[uint32]IinterFace.IConnect),
	}
}

func (this *ConnManager) AddConn(conn IinterFace.IConnect) {
	this.connLock.Lock()
	defer this.connLock.Unlock()

	this.connects[conn.GetTcpConnectID()] = conn

	fmt.Println("[INFO]connection add to ConnManager successfully: conn num = ", this.GetConnNum())

}
func (this *ConnManager) RemoveConn(conn IinterFace.IConnect) {
	this.connLock.Lock()
	defer this.connLock.Unlock()
	delete(this.connects, conn.GetTcpConnectID())
	fmt.Println("connection Remove ConnID=", conn.GetTcpConnectID(), " successfully: conn num = ", this.GetConnNum())

}
func (this *ConnManager) GetConn(connID uint32) (conn IinterFace.IConnect, err error) {
	this.connLock.Lock()
	defer this.connLock.Unlock()
	if conn, ok := this.connects[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connect not found!")
	}
}

func (this *ConnManager) GetConnNum() int {
	//this.connLock.Lock()
	//defer this.connLock.Unlock()
	return len(this.connects)
}

func (this *ConnManager) ClearConn() {
	this.connLock.Lock()
	defer this.connLock.Unlock()

	for connID, conn := range this.connects {
		conn.Stop()
		delete(this.connects, connID)
	}

	fmt.Println("[INFO]Clear All Connections successfully: conn num = ", this.GetConnNum())
}
