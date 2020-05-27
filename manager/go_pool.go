// Author: yann
// Date: 2020/5/25 6:26 下午
// Desc:

package manager

import (
	"fmt"
	"time"
)

// 如果传入任务超时返回的错误信息
var ErrScheduleTimeout = fmt.Errorf("schedule error: timed out")

// 线程池模型
type Pool struct {
	sem  chan struct{} //线程池队列
	work chan func()   //任务队列
}

//***************************************************
//Description : 初始化线程池
//param :       线程最大数量
//param :       任务队列最大数量
//param :       初始化开启多少线程
//return :      线程池
//***************************************************
func NewPool(size, queue, currentSize int) *Pool {
	if currentSize <= 0 && queue > 0 {
		panic("dead queue configuration detected")
	}
	if currentSize > size {
		panic("currentSize > workers")
	}
	p := &Pool{
		sem:  make(chan struct{}, size),
		work: make(chan func(), queue),
	}
	for i := 0; i < currentSize; i++ {
		p.sem <- struct{}{}
		go p.worker(func() {})
	}

	return p
}

//***************************************************
//Description : 传入任务
//param :       任务句柄
//***************************************************
func (p *Pool) Schedule(task func()) {
	p.schedule(task, nil)
}

//***************************************************
//Description : 传入任务设置超时
//param :       超时时间
//param :       任务句柄
//return :      如果任务队列和协程队列均满一直无法写入,超过时间返回超时错误
//***************************************************
func (p *Pool) ScheduleTimeout(timeout time.Duration, task func()) error {
	return p.schedule(task, time.After(timeout))
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrScheduleTimeout
	case p.work <- task: //如果任务队列未满传入任务
		return nil
	case p.sem <- struct{}{}: //如果任务队列已满,创建协程
		go p.worker(task)
		return nil
	}
}

//线程工作
func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()

	task()

	for task := range p.work {
		task()
	}
}
