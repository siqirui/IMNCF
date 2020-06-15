package Ilog

import (
	"testing"
)

func TestStdILoger(t *testing.T) {

	//测试 默认debug输出
	Debug("IMCF debug content1")
	Debug("IMCF debug content2")

	Debugf(" IMCF debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	ResetFlags(BitData | BitLongFile | BitLevel)
	Info("IMCF info content")

	//设置日志前缀，主要标记当前日志模块
	SetPrefix("MODULE")
	Error("IMCF error content")

	//添加标记位
	AddFlag(BitShortFile | BitTime | BitMicroSenconds)
	Stack(" IMCF Stack! ")

	//设置日志写入文件
	SetLogFile("./log", "testfile.log")
	Debug("===> IMCF debug content ~~666")
	Debug("===> IMCF debug content ~~888")
	Error("===> IMCF Error!!!! ~~~555~~~")
	Stack(" IMCF Stack! 111")
	SetLogFile("./log", "testfile1.log")
	Debug("===> IMCF debug content ~~666")
	Debug("===> IMCF debug content ~~888")
	Error("===> IMCF Error!!!! ~~~555~~~")
	Stack(" IMCF Stack! 111")
	//关闭debug调试
	CloseDebug()
	Debug("===> 我不应该出现~！")
	Debug("===> 我不应该出现~！")
	Error("===> IMCF Error  after debug close !!!!")

}
