package xsync_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/dushaoshuai/xsync"
)

func ExampleOnces() {
	onces := xsync.OnceEvery(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10_000_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			onces.Do(func() error {
				fmt.Println("Only once this second:", time.Now().Second())
				return nil
			})
		}()
	}
	wg.Wait()

	// Output:
	// Only once this second: 39
	// Only once this second: 40
	// Only once this second: 41
	// Only once this second: 42
	// Only once this second: 43
}

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
		if i > 1_00_000 {
			t.SkipNow()
		}

		var (
			interval = time.Duration(dNano)
			onces    = xsync.OnceEvery(interval)

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
				go func() {
					for range c {
					}
				}()
			case <-succ:
			}
		}()

		for curr := range c {
			if !prev.IsZero() && prev.Add(interval).After(curr) {
				close(fail)
				t.Fatalf("(%v).Add(%v).After(%v) = %v, want %v", prev, interval, curr, true, false)
			}
			prev = curr
		}
		close(succ)
	})
}
