package TcpServer

import (
	IinterFace "IMCF/IinterFace"
	"IMCF/Utils"
	"fmt"
	"strconv"
)

type MsgHandler struct {
	Handles        map[uint32]IinterFace.IRouter //存放每个MsgEvent 所对应的处理方法的map属性
	WorkerPoolSize uint32                        //任务池数量
	TaskQueue      []chan IinterFace.IRequest    //任务队列
}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Handles:        make(map[uint32]IinterFace.IRouter),
		WorkerPoolSize: Utils.GlobalOBJ.WorkerPoolSize,
		TaskQueue:      make([]chan IinterFace.IRequest, Utils.GlobalOBJ.WorkerPoolSize),
	}
}

//type IMsgHandler interface{
//	DoMsgHandler(request IRequest)			//马上以非阻塞方式处理消息
//	AddRouter(msgId uint32, router IRouter)	//为消息添加具体的处理逻辑
//}

func (this *MsgHandler) SendMsgToTaskQueue(request IinterFace.IRequest) {
	//根据链接ID  轮询分配到没一个任务队列
	workerID := request.GetConnection().GetTcpConnectID() % this.WorkerPoolSize
	fmt.Println("[INFO] add connID = ", request.GetConnection().GetTcpConnectID(), "request msgEvent = ", request.GetMsgEvent(), "to workerID", workerID)
	this.TaskQueue[workerID] <- request
}

func (this *MsgHandler) DoMsgHandler(request IinterFace.IRequest) {
	handle, ok := this.Handles[request.GetMsgEvent()]
	if !ok {
		fmt.Println("api msgEvent = ", request.GetMsgEvent(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	//handle.PreHandle(request)
	handle.Handle(request)
	//handle.PostHandle(request)

}

func (this *MsgHandler) AddRouter(msgEvent uint32, router IinterFace.IRouter) {

	//1 判断当前msgEvent绑定的handle处理方法是否已经存在
	if _, ok := this.Handles[msgEvent]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgEvent)))
	}
	//2 添加msgEvent与Handle的绑定关系
	this.Handles[msgEvent] = router
	fmt.Println("Add Handles msgEvent = ", msgEvent)
}

//启动一个队列
func (this *MsgHandler) StartOneWorker(workID int, taskQueue chan IinterFace.IRequest) {
	fmt.Println("[INFO] worker id = ", workID, "will start")
	for {
		select {
		case request := <-taskQueue:
			this.DoMsgHandler(request)
		}
	}
}

//启动工作池
func (this *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(this.WorkerPoolSize); i++ {
		this.TaskQueue[i] = make(chan IinterFace.IRequest, Utils.GlobalOBJ.MaxWorkerTaskLen)
		go this.StartOneWorker(i, this.TaskQueue[i])
	}
}
