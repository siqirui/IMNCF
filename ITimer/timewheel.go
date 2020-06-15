package ITimer

import (
	"IMCF/Ilog"
	"errors"
	"fmt"
	"sync"
	"time"
)

//时间轮询器为了管理和维护大量的timer调度问题

type TimeWheel struct {
	name         string                    // timewheel名称
	interval     int64                     //刻度的时间间隔ms
	scales       int                       //刻度
	curIndex     int                       //当前时间指针
	maxCap       int                       //每个时间刻度存放的timer的容量
	timerQueue   map[int]map[uint32]*Timer //int轮询刻度，uint32是timerid
	nexTimeWheel *TimeWheel                //
	sync.RWMutex
}

func NewTimeWheel(name string, interval int64, scales int, maxCap int) *TimeWheel {
	timewheel := &TimeWheel{
		name:       name,
		interval:   interval,
		scales:     scales,
		maxCap:     maxCap,
		timerQueue: make(map[int]map[uint32]*Timer, scales),
	}

	for i := 0; i < scales; i++ {
		timewheel.timerQueue[i] = make(map[uint32]*Timer, maxCap)
	}
	Ilog.Info("Init timerWhell name = ", timewheel.name, " is Done!")
	return timewheel
}

func (this *TimeWheel) addTimer(timerID uint32, timer *Timer, forceNext bool) error {
	defer func() error {
		if err := recover(); err != nil {
			errstr := fmt.Sprintf("addTimer function err : %s", err)
			Ilog.Error(errstr)
			return errors.New(errstr)
		}
		return nil
	}()
	delayInterval := timer.unixts - UnixMill()

	if delayInterval >= this.interval {
		delayNum := delayInterval / this.interval

		this.timerQueue[(this.curIndex+int(delayNum))%this.scales][timerID] = timer
	}

	if delayInterval < this.interval && this.nexTimeWheel == nil {
		if forceNext == true {
			this.timerQueue[(this.curIndex+1)%this.scales][timerID] = timer
		} else {
			this.timerQueue[this.curIndex][timerID] = timer
		}
		return nil
	}
	if delayInterval < this.interval {
		return this.nexTimeWheel.AddTimer(timerID, timer)
	}
	return nil
}

func (this *TimeWheel) AddTimer(timerID uint32, timer *Timer) error {
	this.Lock()
	defer this.Unlock()

	return this.addTimer(timerID, timer, false)
}

func (this *TimeWheel) RemoveTimer(timerID uint32) {
	this.Lock()
	defer this.Unlock()

	for i := 0; i < this.scales; i++ {
		if _, ok := this.timerQueue[i][timerID]; ok {
			delete(this.timerQueue[i], timerID)
		}
	}
}

func (this *TimeWheel) AddTimeWheel(next *TimeWheel) {
	this.nexTimeWheel = next
	Ilog.Info("Add timerWhell[", this.name, "]'s next [", next.name, "] is succ!")
}

func (this *TimeWheel) run() {
	for {
		//时间轮每间隔interval一刻度时间，触发转动一次
		time.Sleep(time.Duration(this.interval) * time.Millisecond)

		this.Lock()
		//取出挂载在当前刻度的全部定时器
		curTimers := this.timerQueue[this.curIndex]
		//当前定时器要重新添加 所给当前刻度再重新开辟一个map Timer容器
		this.timerQueue[this.curIndex] = make(map[uint32]*Timer, this.maxCap)
		for tid, timer := range curTimers {
			//这里属于时间轮自动转动，forceNext设置为true
			this.addTimer(tid, timer, true)
		}

		//取出下一个刻度 挂载的全部定时器 进行重新添加 (为了安全起见,待考慮)
		nextTimers := this.timerQueue[(this.curIndex+1)%this.scales]
		this.timerQueue[(this.curIndex+1)%this.scales] = make(map[uint32]*Timer, this.maxCap)
		for tid, timer := range nextTimers {
			this.addTimer(tid, timer, true)
		}

		//当前刻度指针 走一格
		this.curIndex = (this.curIndex + 1) % this.scales

		this.Unlock()
	}
}

func (this *TimeWheel) Run() {
	go this.run()
	Ilog.Info("timerwheel name = ", this.name, " is running...")
}

//获取定时器在一段时间间隔内的Timer
func (this *TimeWheel) GetTimerWithIn(duration time.Duration) map[uint32]*Timer {
	//最终触发定时器的一定是挂载最底层时间轮上的定时器
	//1 找到最底层时间轮
	leaftw := this
	for leaftw.nexTimeWheel != nil {
		leaftw = leaftw.nexTimeWheel
	}

	leaftw.Lock()
	defer leaftw.Unlock()
	//返回的Timer集合
	timerList := make(map[uint32]*Timer)

	now := UnixMill()

	//取出当前时间轮刻度内全部Timer
	for tid, timer := range leaftw.timerQueue[leaftw.curIndex] {
		if timer.unixts-now < int64(duration/1e6) {
			//当前定时器已经超时
			timerList[tid] = timer
			//定时器已经超时被取走，从当前时间轮上 摘除该定时器
			delete(leaftw.timerQueue[leaftw.curIndex], tid)
		}
	}

	return timerList
}
