/**
 * @Author: ChenJunJi
 * @Date: 2023/12/22
 * @Desc:
**/

package timehelper

import (
	"container/heap"
	"fmt"
	"os"
	"time"
)

type Timer struct {
	timeout    time.Time
	interval   time.Duration
	chaseFrame bool
	callback   func()
	cancelled  bool
	seq        uint
}

type TimerHeap struct {
	LogStack func(format string, v ...interface{})
	timers   []*Timer
	seq      uint
}

type TimerOption func(timer *Timer) *Timer

func WithOutChaseFrame() TimerOption {
	return func(timer *Timer) *Timer {
		timer.chaseFrame = false
		return timer
	}
}

func (timer *Timer) Stop() {
	timer.cancelled = true
}

func NewTimerHeap(logStack func(format string, v ...interface{})) *TimerHeap {
	ht := &TimerHeap{seq: 0}
	ht.LogStack = logStack
	heap.Init(ht)
	return ht
}

func (th *TimerHeap) Len() int {
	return len(th.timers)
}

func (th *TimerHeap) Less(i, j int) bool {
	left, right := th.timers[i], th.timers[j]
	if left.timeout.Equal(right.timeout) {
		return left.seq < right.seq
	}

	return left.timeout.Before(right.timeout)
}

func (th *TimerHeap) Swap(i, j int) {
	th.timers[i], th.timers[j] = th.timers[j], th.timers[i]
}

func (th *TimerHeap) Push(x interface{}) {
	th.timers = append(th.timers, x.(*Timer))
}

func (th *TimerHeap) Pop() (ret interface{}) {
	l := len(th.timers)
	th.timers, ret = th.timers[:l-1], th.timers[l-1]
	return
}

func (th *TimerHeap) RunFrame() {
	now := time.Now()
	for {
		if th.Len() <= 0 {
			break
		}
		timer := th.timers[0]
		if timer.cancelled {
			heap.Pop(th)
			continue
		}

		if timer.timeout.After(now) {
			break
		}
		heap.Pop(th)

		th.safeCb(timer)

		if timer.interval > 0 {
			th.seq++
			if !timer.chaseFrame {
				timer.timeout = timer.timeout.Add(timer.interval)
			} else {
				timer.timeout = now.Add(timer.interval)
			}
			timer.seq = th.seq
			heap.Push(th, timer)
		} else {
			timer.cancelled = true
		}
	}
}

func (th *TimerHeap) CancelAll() {
	for _, v := range th.timers {
		v.cancelled = true
	}
}

func (th *TimerHeap) safeCb(timer *Timer) {
	defer func() {
		err := recover()
		if nil == err {
			return
		}
		if nil != th.LogStack {
			th.LogStack("%v", err)
		} else {
			fmt.Fprint(os.Stderr, err)
		}
	}()
	timer.callback()
}

func (th *TimerHeap) SetTimeout(duration time.Duration, cb func()) *Timer {
	if nil == cb {
		return nil
	}
	th.seq++
	timer := new(Timer)
	timer.timeout = time.Now().Add(duration)
	timer.callback = cb
	timer.seq = th.seq
	heap.Push(th, timer)
	return timer
}

func (th *TimerHeap) SetInterval(duration time.Duration, cb func(), ops ...TimerOption) *Timer {
	if nil == cb {
		return nil
	}
	th.seq++
	timer := new(Timer)
	timer.timeout = time.Now().Add(duration)
	timer.callback = cb
	timer.interval = duration
	timer.seq = th.seq
	timer.chaseFrame = true

	for _, op := range ops {
		op(timer)
	}
	heap.Push(th, timer)
	return timer
}
