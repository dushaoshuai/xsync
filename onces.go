package xsync

import (
	"sync"
	"time"
)

// A Onces will perform exactly one successful action at intervals.
// A Onces may be used by multiple goroutines simultaneously.
type Onces struct {
	last     time.Time // the last time one successful action was performed
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
func (t *Onces) needDo() bool {
	return t.last.IsZero() || t.last.Add(t.interval).Before(time.Now())
}

// Do calls the function f if and only if f has never been called successfully
// within this interval. In other words, within current interval,
// f will be invoked each time Do is called unless the previous call to f returns
// without error. After a successful call to f returns, next interval starts.
func (t *Onces) Do(f func() error) error {
	if !t.needDo() {
		return nil
	}

	t.m.Lock()
	defer t.m.Unlock()
	if !t.needDo() {
		return nil
	}
	err := f()
	if err == nil {
		t.last = time.Now()
	}
	return err
}
