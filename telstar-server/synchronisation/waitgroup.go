package synchronisation

import (
	"sync"
	"sync/atomic"
)

type WaitGroupWithCount struct {
	sync.WaitGroup
	count int64
}

func (wg *WaitGroupWithCount) Add(delta int) {
	atomic.AddInt64(&wg.count, int64(delta))
	wg.WaitGroup.Add(delta)
}

func (wg *WaitGroupWithCount) Done() {
	atomic.AddInt64(&wg.count, -1)
	wg.WaitGroup.Done()
}

func (wg *WaitGroupWithCount) GetCount() int {
	return int(atomic.LoadInt64(&wg.count))
}

// Wait() promoted from the embedded field
