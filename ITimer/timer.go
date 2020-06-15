package ITimer

import (
	"time"
)

const (
	HOUR_NAME     = "HOUR"
	HOUR_INTERVAL = 60 * 60 * 1e3
	HOUR_SCALES   = 12

	MINUTE_NAME     = "MINUTE"
	MINUTE_INTERVAL = 60 * 1e3
	MINUTE_SCALES   = 60

	SECOND_NAME     = "SECOND"
	SECOND_INTERVAL = 1e3
	SECOND_SCALES   = 60

	TIMERS_MAX_CAP = 1024
)

/*
    有关时间的几个换算
   	time.Second(秒) = time.Millisecond * 1e3
	time.Millisecond(毫秒) = time.Microsecond * 1e3
	time.Microsecond(微秒) = time.Nanosecond * 1e3

	time.Now().UnixNano() ==> time.Nanosecond (纳秒)
*/

type Timer struct {
	delayFunc *DelayFunc
	unixts    int64
}

//1970-1-1至今 毫秒
func UnixMill() int64 {
	return time.Now().UnixNano() / 1e6
}

//创建一个定时器  指定时间触发的定时器
func NewTimerAt(delayfunc *DelayFunc, unixNano int64) *Timer {
	return &Timer{
		delayFunc: delayfunc,
		unixts:    unixNano / 1e6, //将纳秒转换成对应的毫秒 ms ，定时器以ms为最小精度
	}
}

//创建一个定时器，从当前时间延时duration之后触发
func NewTimerAfter(delayfunc *DelayFunc, duration time.Duration) *Timer {
	return NewTimerAt(delayfunc, time.Now().UnixNano()+int64(duration))
}

//定时启动器
func (this *Timer) Run() {
	go func() {
		now := UnixMill()
		if this.unixts > now {
			time.Sleep(time.Duration(this.unixts-now) * time.Millisecond)
		}
		this.delayFunc.Call()
	}()
}
