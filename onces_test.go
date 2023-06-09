package xsync_test

import (
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/dushaoshuai/xsync"
)

func FuzzDo(f *testing.F) {
	defaultI := 1_00_000
	f.Add(time.Nanosecond.Nanoseconds(), defaultI)
	f.Add(time.Microsecond.Nanoseconds(), defaultI)
	f.Add(time.Millisecond.Nanoseconds(), defaultI)
	f.Add(time.Second.Nanoseconds(), defaultI)
	f.Add(3*time.Second.Nanoseconds(), defaultI)

	f.Fuzz(func(t *testing.T, dNano int64, i int) {
		if dNano > 5*time.Second.Nanoseconds() {
			t.SkipNow()
		}
		if dNano <= 0 {
			t.SkipNow()
		}
		if i > 1_00_000 {
			t.SkipNow()
		}
		if i <= 0 {
			t.SkipNow()
		}

		var (
			interval = time.Duration(dNano)
			onces    = xsync.OnceInterval(interval)

			c    = make(chan time.Time)
			prev time.Time
			fail = make(chan struct{})
			succ = make(chan struct{})

			eg errgroup.Group
		)

		go func() {
		L:
			for j := 0; j < i; j++ {
				select {
				case <-fail:
					break L
				default:
				}
				eg.Go(func() error {
					return onces.Do(func() error {
						c <- time.Now()
						return nil
					})
				})
			}

			eg.Wait()
			close(c)
		}()

		go func() {
			select {
			case <-fail:
				for range c {
				}
			case <-succ:
			}
		}()

		for curr := range c {
			if !prev.IsZero() && prev.Add(interval).After(curr) {
				close(fail)
				eg.Wait()
				t.Fatalf("(%v).Add(%v).After(%v) = %v, want %v", prev, interval, curr, true, false)
			}
			prev = curr
		}
		close(succ)
	})
}
