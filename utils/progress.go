package utils

import (
	"sync"
)

type Progress struct {
	sync.Mutex
	Value    int
	Max      int
	OnStart  bool
	onFinish bool
	ch       chan int
}

func ProgressInit(i int) *Progress {
	ps := &Progress{Max: i, Value: 0, onFinish: false, ch: make(chan int, 10)}
	go ps.workers()
	return ps
}

func (ps *Progress) workers() {
	for value := range ps.ch {
		ps.Lock()
		if ps.Value+value >= ps.Max {
			ps.Value = ps.Max
			ps.onFinish = true
			ps.OnStart = false
		} else {
			ps.Value += value
		}
		ps.Unlock()
	}
}

func (ps *Progress) Add(i int) {
	if !ps.onFinish {
		ps.Lock()
		defer ps.Unlock()

		ps.OnStart = true
		ps.ch <- i
	} else {
		close(ps.ch)
	}
}

func (ps *Progress) Finish() {
	if !ps.onFinish {
		ps.Lock()
		defer ps.Unlock()

		ps.OnStart = true
		ps.ch <- 100
	} else {
		close(ps.ch)
	}
}

func (ps *Progress) Clear() {
	if ps.onFinish {
		ps.Value = 0
		ps.Max = 100
		ps.onFinish = false
		ps.ch = make(chan int, 10)
	}
}
