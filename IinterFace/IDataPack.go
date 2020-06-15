package IinterFace

type IDataPack interface {
	GetHeaderLen() uint32                         //获得包头长度
	Pack(msg IMessage) ([]byte, error)            //封包
	UnPack(data []byte) (msg IMessage, err error) //解包
}
