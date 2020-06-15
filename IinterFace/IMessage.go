package IinterFace

type IMessage interface {
	GetDataLen() uint32  //获取消息数据段长度
	GetMsgEvent() uint32 //获取消息ID
	GetFileID() [8]byte  //获取文件ID
	GetData() []byte     //获取消息内容

	SetMsgEvent(uint32)       //设计消息ID
	SetData([]byte)           //设计消息内容
	SetFileID(fileid [8]byte) //获取文件ID
	SetDataLen(uint32)        //设置消息数据段长度
}
