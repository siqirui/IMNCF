package TcpServer

import (
	IinterFace "IMCF/IinterFace"
	Utils "IMCF/Utils"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type DataPack struct{}

//type IDataPack interface{
//	GetHeaderLen()uint32					//获得包头长度
//	Pack(msg IMessage)([]byte,error)		//封包
//	UnPack(data []byte)(msg IMessage,err error)//解包
//}

//New
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dataPack *DataPack) GetHeadLen() uint32 {

	return uint32(Utils.DATA_UPDATE_HEAD_SIZE)
}

//封包
func (dataPack *DataPack) Pack(msg IinterFace.IMessage) ([]byte, error) {

	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写Event
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgEvent()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetFileID()); err != nil {
		return nil, err
	}
	////写dataLen
	//if msg.GetMsgEvent() == Utils.DATA_UPDATE_AFTER || msg.GetMsgEvent() == Utils.DATA_UPDATE_BEFORE{
	//	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetDataLen()); err != nil {
	//		return nil, err
	//	}
	//}W

	//写data数据
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//封包
func (dataPack *DataPack) Pack1(msg IinterFace.IMessage) ([]byte, error) {

	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写Event
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgEvent()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetFileID()); err != nil {
		return nil, err
	}
	//写dataLen
	if msg.GetMsgEvent() == Utils.DATA_UPDATE_BEFORE || msg.GetMsgEvent() == Utils.DATA_UPDATE {
		if err := binary.Write(dataBuff, binary.BigEndian, msg.GetDataLen()); err != nil {
			return nil, err
		}
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//解包
func (dataPack *DataPack) UnPack(data []byte) (IinterFace.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(data)

	//只解压head的信息，得到dataLen和msgID
	msgbuf := &Message{}

	//读MsgEvent
	if err := binary.Read(dataBuff, binary.BigEndian, &msgbuf.Event); err != nil {
		return nil, err
	}

	fmt.Println("recv msg event:", msgbuf.Event)
	fileid := make([]byte, 8)
	if err := binary.Read(dataBuff, binary.BigEndian, fileid); err != nil {
		return nil, err
	}
	copy(msgbuf.FileID[:], fileid)
	str := string(fileid[:])
	fmt.Println(str)
	fmt.Println("recv msg fileid:", str)

	//读dataLen
	//if msgbuf.Event != Utils.DATA_UPDATE_AFTER{
	//	if err := binary.Read(dataBuff, binary.BigEndian, &msgbuf.DataLen); err != nil {
	//		return nil, err
	//	}
	//}else{
	//	msgbuf.DataLen = 0
	//}
	//msgbuf.DataLen = msgbuf.DataLen - 4
	//fmt.Println("recv msgbuf.DataLen:" ,msgbuf.DataLen)

	//判断dataLen的长度是否超出我们允许的最大包长度
	if Utils.GlobalOBJ.MaxPackSize > 0 && msgbuf.DataLen > Utils.GlobalOBJ.MaxPackSize {
		return nil, errors.New("Too large msg data recieved")
	}
	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据

	return msgbuf, nil
}
