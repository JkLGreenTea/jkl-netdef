package waitgroup

import (
	"go.uber.org/atomic"
	"sync"
)

type WaitGroup struct {
	Wg    *sync.WaitGroup // WaitGroup.
	Count *atomic.Uint32  // Count goroutine.
}

// New - new  WaitGroup.
func New() *WaitGroup {
	return &WaitGroup{
		Wg:    &sync.WaitGroup{},
		Count: atomic.NewUint32(0),
	}
}

// Add - sync.WaitGroup.Add()
func (wg *WaitGroup) Add(delta int) {
	wg.Wg.Add(delta)
	wg.Count.Add(uint32(delta))
}

// Done - sync.WaitGroup.Done()
func (wg *WaitGroup) Done() {
	wg.Wg.Done()
	wg.Count.Sub(1)
}

// Wait - sync.WaitGroup.Wait()
func (wg *WaitGroup) Wait() {
	wg.Wg.Wait()
}
