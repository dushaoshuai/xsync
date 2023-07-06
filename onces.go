package xsync

import (
	"sync"
	"sync/atomic"
	"time"
)

// Onces is an object that will try to (if asked) perform a successful action
// only if a specified interval has elapsed since the last successful action.
// An Onces must not be copied after first use.
// An Onces may be used by multiple goroutines simultaneously.
type Onces struct {
	// next is the next time one successful action may be performed,
	// represented in Unix nanoseconds.
	next int64
	// at least this interval between two successful actions
	interval time.Duration

	m sync.Mutex
}

// OnceInterval returns a new Onces.
// It will panic if interval is less than or equal to 0.
func OnceInterval(interval time.Duration) *Onces {
	if interval <= 0 {
		panic("xsync: OnceInterval: interval must be greater than 0")
	}
	return &Onces{
		next:     time.Now().UnixNano(), // the first actions may be performed immediately without having to wait
		interval: interval,
	}
}

// Do calls the function f if and only if the specified
// interval has elapsed since the last successful call to f.
// Do considers a call to f succeed if f returns nil.
// It's ok if f has different values in each invocation.
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
