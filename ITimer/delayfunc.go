package ITimer

import (
	"IMCF/Ilog"
	"fmt"
	"reflect"
)

//   定义一个延迟调用函数  定时器超时时 调用回调

type DelayFunc struct {
	f    func(...interface{}) //延迟函数调用原型
	args []interface{}        //延迟调用函数传递的形参
}

func NewDelayFunc(f func(v ...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f:    f,
		args: args,
	}
}

//打印当前延迟函数的信息，用于日志记录
func (this *DelayFunc) String() string {
	return fmt.Sprintf("{DelayFun:%s, args:%v}", reflect.TypeOf(this.f).Name(), this.args)
}

//执行延迟函数---如果执行失败，抛出异常
func (this *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			Ilog.Error(this.String(), "Call err: ", err)
		}
	}()

	//调用定时器超时函数
	this.f(this.args...)
}
