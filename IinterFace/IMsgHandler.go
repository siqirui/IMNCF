package IinterFace

//消息管理器
type IMsgHandler interface {
	DoMsgHandler(request IRequest)             //马上以非阻塞方式处理消息
	AddRouter(msgEvent uint32, router IRouter) //为消息添加具体的处理逻辑
	StartWorkerPool()                          //启动任务池
	SendMsgToTaskQueue(request IRequest)       //发送消息进入任务队列
}
