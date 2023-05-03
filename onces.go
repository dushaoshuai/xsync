package xsync

import (
	"sync"
	"sync/atomic"
	"time"
)

// A Onces will perform exactly one successful action at intervals.
// A Onces may be used by multiple goroutines simultaneously.
type Onces struct {
	// next is the next time one successful action may be performed,
	// represented in Unix nanoseconds.
	next int64

	interval time.Duration
	m        sync.Mutex
}

// OnceEvery returns a Onces that will perform exactly one successful action within each interval.
func OnceEvery(interval time.Duration) *Onces {
	return &Onces{
		next:     time.Now().UnixNano(),
		interval: interval,
	}
}

// Do calls the function f if and only if f has never been called successfully
// within this interval. In other words, within current interval,
// f will be invoked each time Do is called unless the previous call to f returns
// without error. After a successful call to f returns, next interval starts.
func (o *Onces) Do(f func() error) error {
	if atomic.LoadInt64(&o.next) > time.Now().UnixNano() {
		return nil
	}

	o.m.Lock()
	defer o.m.Unlock()
	if o.next > time.Now().UnixNano() {
		return nil
	}
	err := f()
	if err == nil {
		atomic.StoreInt64(&o.next, time.Now().Add(o.interval).UnixNano())
	}
	return err
}
