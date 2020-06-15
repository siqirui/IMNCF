package Worker

import (
	"IMCF/IinterFace"
	"IMCF/Inet/TcpServer"
	"IMCF/Utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"logging"
	"protobuffer"
	"service/impl"
	"strconv"
	"time"
	"utils"
	"web"
)

const (
	FILEIDLEN = 8
	TOKENLEN  = 32
)

var ConnectFiles = make(map[uint32]*RecvFile)

type RecvFileBefore struct {
	Event       uint32
	FileID      [FILEIDLEN]byte
	Len         int32
	ShardingNum byte
	UserID      uint32
	Token       [TOKENLEN]byte
	FileType    byte
	FileName    string
}
type RecvFileAfter struct {
	Event  uint32
	FileID [FILEIDLEN]byte
}

type DoRecvFile struct {
	ShardingSeq byte
	FileDataMap map[int][]byte
}

type RecvFile struct {
	Before RecvFileBefore
	After  RecvFileAfter
	Recv   DoRecvFile
}

type Router struct {
	TcpServer.BaseRouter
}

func NewRecvFileService() *RecvFile {
	return &RecvFile{
		RecvFileBefore{},
		RecvFileAfter{},
		DoRecvFile{},
	}
}
func (this *RecvFile) unpackMessageDo(data []byte) error {

	dataBytes := bytes.NewBuffer(data)
	if err := binary.Read(dataBytes, binary.BigEndian, &this.Recv.ShardingSeq); err != nil {
		logging.Logger.Debugf("unpackMessageDo-->read ShardingSeq error data: %x,ShardingSeq : %x", dataBytes.Bytes(), this.Recv.ShardingSeq)
		return err
	}
	len := len(data) - 1
	mdata := make([]byte, len)
	if err := binary.Read(dataBytes, binary.LittleEndian, mdata); err != nil {
		logging.Logger.Debugf("unpackMessageDo-->read fileinfo err fileinfo:%x databytes: %x", mdata, dataBytes.Bytes())

	}
	Seq := int(this.Recv.ShardingSeq)
	logging.Logger.Debugf("unpackMessageDo-->read seqNum : %x", this.Recv.ShardingSeq)
	logging.Logger.Debug("unpackMessageDo-->read seqNum :", Seq)
	this.Recv.FileDataMap[Seq-1] = mdata

	return nil
}
func (this *RecvFile) unpackMessageBefore(data []byte) error {

	dataBytes := bytes.NewBuffer(data)
	if err := binary.Read(dataBytes, binary.BigEndian, &this.Before.ShardingNum); err != nil {
		logging.Logger.Debugf("unpackMessageBefore-->read ShardingNum err ShardingNum:%x databytes: %x", this.Before.ShardingNum, dataBytes.Bytes())
		return err
	}
	if err := binary.Read(dataBytes, binary.BigEndian, &this.Before.UserID); err != nil {
		logging.Logger.Debugf("unpackMessageBefore-->read UserID err UserID:%d databytes: %x", this.Before.UserID, dataBytes.Bytes())
		return err
	}
	//token := make([]byte,TOKENLEN)
	token := make([]byte, TOKENLEN)
	if err := binary.Read(dataBytes, binary.BigEndian, token); err != nil {
		logging.Logger.Debugf("unpackMessageBefore-->read token err token:%x databytes: %x", token, dataBytes.Bytes())
		return err
	}
	copy(this.Before.Token[:7], token)
	logging.Logger.Debugf("unpackMessageBefore-->read this.Before.Token :%x databytes: %x", this.Before.Token, dataBytes.Bytes())

	if err := binary.Read(dataBytes, binary.BigEndian, &this.Before.FileType); err != nil {
		logging.Logger.Debugf("unpackMessageBefore-->read FileType err FileType:%x databytes: %x", this.Before.FileType, dataBytes.Bytes())
		return err
	}

	return nil
}
func (this *RecvFile) packMessageBefore(code byte) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.BigEndian, code); err != nil {
		logging.Logger.Debugf("packMessageBefore-->write code err code:%x databytes: %x", code, dataBuff.Bytes())
		return nil, err
		//TODO 返回错误 log记录
	}

	return dataBuff.Bytes(), nil

}
func (this *RecvFile) unpackMessageAfter(data []byte) error {

	return nil
}
func (this *RecvFile) packMessageAfter(errlist []byte, url1 []byte) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if len(errlist) > 0 {
		fmt.Println("在pack中存在丢包", errlist)
		if err := binary.Write(dataBuff, binary.BigEndian, len(errlist)); err != nil {
			logging.Logger.Debugf("packMessageAfter-->write errlistlen err errlistlen:%d databytes: %x", len(errlist), dataBuff.Bytes())
			return nil, err
		}
		if err := binary.Write(dataBuff, binary.BigEndian, errlist); err != nil {
			logging.Logger.Debugf("packMessageAfter-->write errlist err errlist:%x databytes: %x", errlist, dataBuff.Bytes())
			return nil, err
		}
		logging.Logger.Debugf("packMessageAfter pack %x", dataBuff.Bytes())
		return dataBuff.Bytes(), nil
	} else if len(url1) <= 0 {

		if err := binary.Write(dataBuff, binary.BigEndian, len(url1)); err != nil {
			logging.Logger.Debugf("packMessageAfter-->write errcode  err errcode:%d databytes: %x", len(url1), dataBuff.Bytes())
			return nil, err
		}
		if err := binary.Write(dataBuff, binary.BigEndian, len(url1)); err != nil {
			logging.Logger.Debugf("packMessageAfter-->write errcode err errcode:%d databytes: %x", len(url1), dataBuff.Bytes())
			return nil, err
		}
		logging.Logger.Debugf("packMessageAfter pack %x", dataBuff.Bytes())
		return dataBuff.Bytes(), nil
	}
	var lostlen byte = 0
	if err := binary.Write(dataBuff, binary.BigEndian, lostlen); err != nil {
		return nil, err
	}
	url1 = []byte("http://kbim-res.oss-cn-beijing.aliyuncs.com/chat_file/2020-05-13/494848.amr")
	url1len := byte(uint32(len(url1)))
	if err := binary.Write(dataBuff, binary.BigEndian, url1len); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.BigEndian, url1); err != nil {
		return nil, err
	}

	logging.Logger.Debugf("packMessageAfter pack %x", dataBuff.Bytes())
	return dataBuff.Bytes(), nil

}
func (this *RecvFile) CheckSuffix(suffixNum int) string {
	switch suffixNum {
	case 0:
		return "amr"
	case 1:
		return "amr"
	}
	return "amr"
}

func (this *RecvFile) HandleBefore(request IinterFace.IRequest) {
	message := request.GetData()
	if err := this.unpackMessageBefore(message.GetData()); err != nil {

	}

	var sdktype int32 = 0
	if err := this.CheckToken(int32(this.Before.UserID), sdktype, string(this.Before.Token[:]), string(this.Before.FileID[:])); !err {
		logging.Logger.Debugf("check token err --> userid:%d  token:%s ", this.Before.UserID, string(this.Before.Token[:]))
		//this.Before.FileName = fmt.Sprintf("./%s.%s",Str,this.CheckSuffix(int(this.Before.FileType)))
		this.Before.FileID = message.GetFileID()
		this.Before.Event = request.GetMsgEvent()
		//this.Before.Len = int32(message.GetDataLen())
		response, err := this.packMessageBefore(0x02)
		if err != nil {
			//TODO 容错处理
		}
		if err := request.GetConnection().SendMsg(this.Before.Event, this.Before.FileID, response); err != nil {
			//TODO 容错处理
		}
	}
	test := message.GetFileID()
	var Str string
	for i := 0; i < 8; i++ {
		if test[i] != 0 {
			var a int = int(test[i])
			b := strconv.Itoa(a)
			Str = Str + b
		}
	}
	//str:=string(test[:])

	this.Before.FileName = fmt.Sprintf("./%s.%s", Str, this.CheckSuffix(int(this.Before.FileType)))
	this.Before.FileID = message.GetFileID()
	this.Before.Event = request.GetMsgEvent()
	this.Before.Len = int32(message.GetDataLen())
	response, err := this.packMessageBefore(0x01)
	if err != nil {
		//TODO 容错处理
	}
	if err := request.GetConnection().SendMsg(this.Before.Event, this.Before.FileID, response); err != nil {
		//TODO 容错处理
	}

	this.Recv.FileDataMap = make(map[int][]byte)
}

func (this *RecvFile) HandleAfter(request IinterFace.IRequest) {

	testlist := make([]byte, 1024)
	var j int = 0
	this.After.Event = request.GetMsgEvent()
	this.After.FileID = request.GetData().GetFileID()

	for i := 0; i < int(this.Before.ShardingNum); i++ {
		if this.Recv.FileDataMap[i] == nil {
			testlist[j] = byte(i)
			j++
		}
	}
	lostList := testlist[:j]
	if j != 0 {
		//TODO response errdata

		url := []byte{0}
		response, err := this.packMessageAfter(lostList, url)
		if err != nil {
			fmt.Println("11111111111")
		}
		fmt.Println("存在丢包！", lostList)
		request.GetConnection().SendMsg(this.After.Event, this.After.FileID, response)
		return
	}
	var fileData []byte
	for i := 0; i < int(this.Before.ShardingNum); i++ {
		fileData = append(fileData, this.Recv.FileDataMap[i]...)
	}

	_, remoteUrl, err := impl.NewOssService().Upload(int32(this.Before.UserID), utils.FormatTimeyyyyMMdd(time.Now()), string(this.Before.FileID[:]), this.Before.FileType, fileData)

	if err != nil {
		//TODO 应该响应失败，但是接口暂时未定义该错误码，协商后再定

		url := []byte{0}
		lostList := []byte{0}
		response, err := this.packMessageAfter(lostList, url)
		if err != nil {
			//TODO 暂不处理
		}
		logging.Logger.Debug("upload error ", err)
		request.GetConnection().SendMsg(this.After.Event, this.After.FileID, response)
		return
	}
	response, err := this.packMessageAfter(lostList, []byte(remoteUrl))
	request.GetConnection().SendMsg(this.After.Event, this.After.FileID, response)

	//var fd *os.File
	//var err1 error
	//if _ , err := os.Stat(this.Before.FileName);!os.IsNotExist(err) {
	//	//文件存在，打开
	//	fd, err1 = os.OpenFile(this.Before.FileName, os.O_APPEND|os.O_RDWR, 0644)
	//	fmt.Println(err1)
	//} else {
	//	//文件不存在，创建
	//	fd, err1 = os.OpenFile(this.Before.FileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	//}
	//
	//defer fd.Close()
	//
	//fd.Write(fileData)

	delete(ConnectFiles, request.GetConnection().GetTcpConnectID())
}

func (this *RecvFile) HandleRecv(request IinterFace.IRequest) {
	message := request.GetData()
	if err := this.unpackMessageDo(message.GetData()); err != nil {

	}

}
func (this *Router) Handle(request IinterFace.IRequest) {

	var handle *RecvFile = nil
	handle = ConnectFiles[request.GetConnection().GetTcpConnectID()]
	logging.Logger.Debugf("connectid : ", request.GetConnection().GetTcpConnectID())
	if handle == nil {
		handle = NewRecvFileService()
		ConnectFiles[request.GetConnection().GetTcpConnectID()] = handle
	}

	if request.GetData().GetDataLen() <= 0 {

	}

	switch request.GetMsgEvent() {
	case Utils.DATA_UPDATE_BEFORE:
		fmt.Println(".....")
		handle.HandleBefore(request)
	case Utils.DATA_UPDATE:
		fmt.Println("......")
		handle.HandleRecv(request)
	case Utils.DATA_UPDATE_AFTER:
		fmt.Println("......")
		handle.HandleAfter(request)
	default:
		fmt.Println("异常数据")
	}
}

func (this *RecvFile) CheckToken(userid int32, sdktype int32, fileid string, token string) bool {

	req := &protobuffer.RequestToken{UserId: userid, TType: web.ConvertTerminalType(sdktype), RequestId: fileid}

	return impl.NewLoginService().CheckToken(req, token)
}
