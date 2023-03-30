package xsync

import (
	"sync"
	"sync/atomic"
	"time"
)

// A Onces will perform exactly one successful action at intervals.
// A Onces may be used by multiple goroutines simultaneously.
type Onces struct {
	// last is the last time one successful action was performed,
	// represented in Unix nanoseconds.
	last     int64 // todo use atomic.Value
	interval time.Duration
	m        sync.Mutex
}

// OnceEvery returns a Onces that will perform exactly one successful action within each interval.
func OnceEvery(interval time.Duration) *Onces {
	return &Onces{
		interval: interval,
	}
}

// needDo reports whether an action should be performed.
func (o *Onces) needDo(atomically bool) bool {
	var last int64
	if atomically {
		last = atomic.LoadInt64(&o.last)
	} else {
		last = o.last
	}
	return last == 0 || time.Unix(0, last).Add(o.interval).Before(time.Now())
}

// Do calls the function f if and only if f has never been called successfully
// within this interval. In other words, within current interval,
// f will be invoked each time Do is called unless the previous call to f returns
// without error. After a successful call to f returns, next interval starts.
func (o *Onces) Do(f func() error) error {
	if !o.needDo(true) {
		return nil
	}

	o.m.Lock()
	defer o.m.Unlock()
	if !o.needDo(false) {
		return nil
	}
	err := f()
	if err == nil {
		atomic.StoreInt64(&o.last, time.Now().UnixNano())
	}
	return err
}
