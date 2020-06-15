package ITimer

import (
	"IMCF/Ilog"
	"math"

	//Iinterface "IMCF/IinterFace"
	_ "fmt"
	"sync"
	"time"
)

const (
	//默认缓冲触发函数队列大小
	MAX_CHAN_BUFF = 2048
	//默认最大误差时间
	MAX_TIME_DELAY = 100
)

type TimerScheduler struct {
	timewheel   *TimeWheel
	idGen       uint32          //定时器编号累加
	triggerChan chan *DelayFunc //已触发定时器的channel
	sync.RWMutex
}

func NewTimerScheduler() *TimerScheduler {
	secondTW := NewTimeWheel(SECOND_NAME, SECOND_INTERVAL, SECOND_SCALES, TIMERS_MAX_CAP)
	minuteTW := NewTimeWheel(MINUTE_NAME, MINUTE_INTERVAL, MINUTE_SCALES, TIMERS_MAX_CAP)
	hourTW := NewTimeWheel(HOUR_NAME, HOUR_INTERVAL, HOUR_SCALES, TIMERS_MAX_CAP)

	hourTW.AddTimeWheel(minuteTW)
	minuteTW.AddTimeWheel(secondTW)

	secondTW.Run()
	minuteTW.Run()
	hourTW.Run()

	return &TimerScheduler{
		timewheel:   hourTW,
		triggerChan: make(chan *DelayFunc, MAX_CHAN_BUFF),
	}
}

func (this *TimerScheduler) CreateTimerAt(delayfunc *DelayFunc, unixNano int64) (uint32, error) {
	this.Lock()
	defer this.Unlock()

	this.idGen++
	return this.idGen, this.timewheel.AddTimer(this.idGen, NewTimerAt(delayfunc, unixNano))
}

func (this *TimerScheduler) CreateTimerAfter(delayfunc *DelayFunc, duration time.Duration) (uint32, error) {
	this.Lock()
	defer this.Unlock()

	this.idGen++
	return this.idGen, this.timewheel.AddTimer(this.idGen, NewTimerAfter(delayfunc, duration))
}

func (this *TimerScheduler) CancelTimer(timerID uint32) {
	this.Lock()
	defer this.Unlock()

	this.timewheel.RemoveTimer(timerID)
}

func (this *TimerScheduler) GetTriggerChan() chan *DelayFunc {
	return this.triggerChan
}

//非阻塞的方式启动timer调度器
func (this *TimerScheduler) Start() {
	go func() {
		for {
			now := UnixMill()
			//获取最近MAX_TIME_DELAY 毫秒的超时的所有定时器
			timerList := this.timewheel.GetTimerWithIn(MAX_TIME_DELAY * time.Millisecond)

			for _, timer := range timerList {
				if math.Abs(float64(now-timer.unixts)) > MAX_TIME_DELAY {
					//告警
					Ilog.Error("wang call at ", timer.unixts, ";real call at", now, ";delay ", now-timer.unixts)
				}
				this.triggerChan <- timer.delayFunc
			}
			//有待优化
			time.Sleep(MAX_TIME_DELAY / 2 * time.Millisecond)
		}

	}()
}

//时间轮定器 自动调度
func NewAutoExecTimerScheduler() *TimerScheduler {
	//创建调度器
	autoExecScheduler := NewTimerScheduler()
	//启动调度器
	autoExecScheduler.Start()

	go func() {
		delayFuncChan := autoExecScheduler.GetTriggerChan()
		for delayfunc := range delayFuncChan {
			go delayfunc.Call()
		}

	}()

	return autoExecScheduler
}
