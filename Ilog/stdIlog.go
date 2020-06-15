package Ilog

import "os"

/*
   全局默认提供一个Log对外句柄，可以直接使用API系列调用
   全局日志对象 StdZinxLog
*/

var StdILog = NewIMCFLog(os.Stderr, "", BitDefault)

//获取IMCFLogger 标记位
func Flags() int {
	return StdILog.Flags()
}

//设置IMCFLogger标记位
func ResetFlags(flag int) {
	StdILog.ResetFlags(flag)
}

//添加flag标记
func AddFlag(flag int) {
	StdILog.AddFlag(flag)
}

//设置IMCFLogger 日志头前缀
func SetPrefix(prefix string) {
	StdILog.SetPrefix(prefix)
}

//设置IMCFLogger绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	StdILog.SetLogFile(fileDir, fileName)
}

//设置关闭debug
func CloseDebug() {
	StdILog.CloseDebug()
}

//设置打开debug
func OpenDebug() {
	StdILog.OpenDebug()
}

// Debug
func Debugf(format string, v ...interface{}) {
	StdILog.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	StdILog.Debug(v...)
}

//Info
func Infof(format string, v ...interface{}) {
	StdILog.Infof(format, v...)
}

func Info(v ...interface{}) {
	StdILog.Info(v...)
}

// Warn
func Warnf(format string, v ...interface{}) {
	StdILog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdILog.Warn(v...)
}

// Error
func Errorf(format string, v ...interface{}) {
	StdILog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdILog.Error(v...)
}

// Fatal 需要终止程序
func Fatalf(format string, v ...interface{}) {
	StdILog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdILog.Fatal(v...)
}

// Panic
func Panicf(format string, v ...interface{}) {
	StdILog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdILog.Panic(v...)
}

// Stack
func Stack(v ...interface{}) {
	StdILog.Stack(v...)
}

func init() {
	//因为StdILog对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	//一般的StdILog对象 calldDepth=2, StdILog的calldDepth=3
	StdILog.calldDepth = 3
}
