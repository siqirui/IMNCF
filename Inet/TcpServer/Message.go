package TcpServer

type Message struct {
	Event   uint32 //消息类型
	FileID  [8]byte
	DataLen uint32 //消息长度
	Data    []byte //消息体
	//IsEncrypt	bool //是否加密
	//ReqOrResp	bool //请求or响应   true为Request false为Response
}

//type IMessage interface {
//	GetDataLen() uint32	//获取消息数据段长度
//	GetMsgEvent() uint32	//获取消息ID
//	GetData() []byte	//获取消息内容
//
//	SetMsgId(uint32)	//设计消息ID
//	SetData([]byte)		//设计消息内容
//	SetDataLen(uint32)	//设置消息数据段长度
//}

//创建一个Message消息包
func NewMsgPackage(event uint32, fileid [8]byte, data []byte) *Message {
	return &Message{
		Event:   event,
		FileID:  fileid,
		DataLen: uint32(int32(len(data))),
		Data:    data,
	}
}

func (this *Message) GetDataLen() uint32 {
	return this.DataLen
}

func (this *Message) GetMsgEvent() uint32 {
	return this.Event
}

func (this *Message) GetData() []byte {
	return this.Data
}

func (this *Message) GetFileID() [8]byte {
	return this.FileID
}

func (this *Message) SetMsgEvent(event uint32) {
	this.Event = event
}

func (this *Message) SetFileID(fileid [8]byte) {
	this.FileID = fileid
}

func (this *Message) SetDataLen(len uint32) {
	this.DataLen = len
}

func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
