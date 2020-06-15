package Ilog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	LOG_MAX_BUF = 1024 * 1024
)

const (
	BitData          = 1 << iota                            //日期标志
	BitTime                                                 //时间标志
	BitMicroSenconds                                        //毫秒级时间标志
	BitLongFile                                             //完整文件路径名
	BitShortFile                                            //文件名
	BitLevel                                                //当前日志级别
	BitStdFlag       = BitData | BitTime                    //标准头部日志格式
	BitDefault       = BitLevel | BitShortFile | BitStdFlag //默认日志头部格式
)

const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogError
	LogPanic
	LogFatal
)

var levels = []string{
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

type IMCFLogger struct {
	Mutex      sync.Mutex   //确保多协程读写文件，防止文件内容混乱，做到协程安全
	prefix     string       //每行log日志的前缀字符串,拥有日志标记
	flag       int          //日志标记位
	out        io.Writer    //日志输出的文件描述符
	buf        bytes.Buffer //输出的缓冲区
	file       *os.File     //当前日志绑定的输出文件
	debugClose bool         //是否打印调试debug信息
	calldDepth int          //获取日志文件名和代码上述的runtime.Call 的函数调用层数
}

func NewIMCFLog(out io.Writer, prefix string, flag int) *IMCFLogger {
	ilog := &IMCFLogger{
		out:        out,
		prefix:     prefix,
		flag:       flag,
		file:       nil,
		debugClose: false,
		calldDepth: 2,
	}
	runtime.SetFinalizer(ilog, CleanIMCFLog)
	return ilog
}
func CleanIMCFLog(log *IMCFLogger) {

}

func (log *IMCFLogger) formatHeader(buf *bytes.Buffer, timestamp time.Time, file string, line int, level int) {
	//如果当前前缀字符串不为空，那么需要先写前缀
	if log.prefix != "" {
		buf.WriteByte('<')
		buf.WriteString(log.prefix)
		buf.WriteByte('>')
	}

	//已经设置了时间相关的标识位,那么需要加时间信息在日志头部
	if log.flag&(BitData|BitTime|BitMicroSenconds) != 0 {
		//日期位被标记
		if log.flag&BitData != 0 {
			year, month, day := timestamp.Date()
			itoa(buf, year, 4)
			buf.WriteByte('/') // "2019/"
			itoa(buf, int(month), 2)
			buf.WriteByte('/') // "2019/04/"
			itoa(buf, day, 2)
			buf.WriteByte(' ') // "2019/04/11 "
		}

		//时钟位被标记
		if log.flag&(BitTime|BitMicroSenconds) != 0 {
			hour, min, sec := timestamp.Clock()
			itoa(buf, hour, 2)
			buf.WriteByte(':') // "11:"
			itoa(buf, min, 2)
			buf.WriteByte(':') // "11:15:"
			itoa(buf, sec, 2)  // "11:15:33"
			//微秒被标记
			if log.flag&BitMicroSenconds != 0 {
				buf.WriteByte('.')
				itoa(buf, timestamp.Nanosecond()/1e3, 6) // "11:15:33.123123
			}
			buf.WriteByte(' ')
		}

		// 日志级别位被标记
		if log.flag&BitLevel != 0 {
			buf.WriteString(levels[level])
		}

		//日志当前代码调用文件名名称位被标记
		if log.flag&(BitShortFile|BitLongFile) != 0 {
			//短文件名称
			if log.flag&BitShortFile != 0 {
				short := file
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						//找到最后一个'/'之后的文件名称  如:/home/go/src/IMFC.go 得到 "IMFC.go"
						short = file[i+1:]
						break
					}
				}
				file = short
			}
			buf.WriteString(file)
			buf.WriteByte(':')
			itoa(buf, line, -1) //行数
			buf.WriteString(": ")
		}
	}
}

func (log *IMCFLogger) OutPut(level int, s string) error {

	now := time.Now() // 得到当前时间
	var file string   //当前调用日志接口的文件名称
	var line int      //当前代码行数
	log.Mutex.Lock()
	defer log.Mutex.Unlock()

	if log.flag&(BitShortFile|BitLongFile) != 0 {
		log.Mutex.Unlock()
		var ok bool
		//得到当前调用者的文件名称和执行到的代码行数
		_, file, line, ok = runtime.Caller(log.calldDepth)
		if !ok {
			file = "unknown-file"
			line = 0
		}
		log.Mutex.Lock()
	}

	//清零buf
	log.buf.Reset()
	//写日志头
	log.formatHeader(&log.buf, now, file, line, level)
	//写日志内容
	log.buf.WriteString(s)
	//补充回车
	if len(s) > 0 && s[len(s)-1] != '\n' {
		log.buf.WriteByte('\n')
	}

	//将填充好的buf 写到IO输出上
	_, err := log.out.Write(log.buf.Bytes())
	return err
}

//Debug
func (log *IMCFLogger) Debugf(format string, v ...interface{}) {
	if log.debugClose == true {
		return
	}
	_ = log.OutPut(LogDebug, fmt.Sprintf(format, v...))
}

func (log *IMCFLogger) Debug(v ...interface{}) {
	if log.debugClose == true {
		return
	}
	_ = log.OutPut(LogDebug, fmt.Sprintln(v...))
}

//Info
func (log *IMCFLogger) Infof(format string, v ...interface{}) {
	_ = log.OutPut(LogInfo, fmt.Sprintf(format, v...))
}

func (log *IMCFLogger) Info(v ...interface{}) {
	_ = log.OutPut(LogInfo, fmt.Sprintln(v...))
}

//Warn
func (log *IMCFLogger) Warnf(format string, v ...interface{}) {
	_ = log.OutPut(LogWarn, fmt.Sprintf(format, v...))
}

func (log *IMCFLogger) Warn(v ...interface{}) {
	_ = log.OutPut(LogWarn, fmt.Sprintln(v...))
}

//Error

func (log *IMCFLogger) Errorf(format string, v ...interface{}) {
	_ = log.OutPut(LogError, fmt.Sprintf(format, v...))
}

func (log *IMCFLogger) Error(v ...interface{}) {
	_ = log.OutPut(LogError, fmt.Sprintln(v...))
}

//Fatal
func (log *IMCFLogger) Fatalf(format string, v ...interface{}) {
	_ = log.OutPut(LogFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (log *IMCFLogger) Fatal(v ...interface{}) {
	_ = log.OutPut(LogFatal, fmt.Sprintln(v...))
	os.Exit(1)
}

//Panic
func (log *IMCFLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	_ = log.OutPut(LogPanic, fmt.Sprintf(format, s))
	panic(s)
}

func (log *IMCFLogger) Panic(v ...interface{}) {
	s := fmt.Sprintln(v...)
	_ = log.OutPut(LogPanic, s)
	panic(s)
}

//Stack
func (log *IMCFLogger) Stack(v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, LOG_MAX_BUF)
	n := runtime.Stack(buf, true) //得到当前堆栈信息
	s += string(buf[:n])
	s += "\n"
	_ = log.OutPut(LogError, s)
}

//Flags
func (log *IMCFLogger) Flags() int {
	log.Mutex.Lock()
	defer log.Mutex.Unlock()
	return log.flag
}

//重置日志flags
func (log *IMCFLogger) ResetFlags(flag int) {
	log.Mutex.Lock()
	defer log.Mutex.Unlock()
	log.flag = flag
}

//添加标记
func (log *IMCFLogger) AddFlag(flag int) {
	log.Mutex.Lock()
	defer log.Mutex.Unlock()
	log.flag |= flag
}

//设置日志的 用户自定义前缀字符串
func (log *IMCFLogger) SetPrefix(prefix string) {
	log.Mutex.Lock()
	defer log.Mutex.Unlock()
	log.prefix = prefix
}

//设置日志文件输出
func (log *IMCFLogger) SetLogFile(fileDir string, fileName string) {
	var file *os.File

	//创建日志文件夹
	_ = mkdirLog(fileDir)

	fullPath := fileDir + "/" + fileName
	if log.checkFileExist(fullPath) {
		//文件存在，打开
		file, _ = os.OpenFile(fullPath, os.O_APPEND|os.O_RDWR, 0644)
	} else {
		//文件不存在，创建
		file, _ = os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	}

	log.Mutex.Lock()
	defer log.Mutex.Unlock()

	//关闭之前绑定的文件
	log.closeFile()
	log.file = file
	log.out = file
}

//关闭日志绑定的文件
func (log *IMCFLogger) closeFile() {
	if log.file != nil {
		_ = log.file.Close()
		log.file = nil
		log.out = os.Stderr
	}
}

func (log *IMCFLogger) CloseDebug() {
	log.debugClose = true
}

func (log *IMCFLogger) OpenDebug() {
	log.debugClose = false
}

//判断日志文件是否存在
func (log *IMCFLogger) checkFileExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func mkdirLog(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0775); err != nil {
			if os.IsPermission(err) {
				e = err
			}
		}
	}
	return
}

func itoa(buf *bytes.Buffer, i int, wid int) {
	var u uint = uint(i)
	if u == 0 && wid <= 1 {
		buf.WriteByte('0')
		return
	}

	var buffer [32]byte
	bufPoint := len(buffer)
	for ; u > 0 || wid > 0; u /= 10 {
		bufPoint--
		wid--
		buffer[bufPoint] = byte(u%10) + '0'
	}

	for bufPoint < len(buffer) {
		buf.WriteByte(buffer[bufPoint])
		bufPoint++
	}
}
