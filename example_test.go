package xsync_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/dushaoshuai/xsync"
)

func ExampleOnces() {
	onces := xsync.OnceInterval(time.Second)

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

	// Sample Output:
	// Only once this second: 30
	// Only once this second: 31
	// Only once this second: 32
	// Only once this second: 33
	// Only once this second: 34
	// Only once this second: 35
	// Only once this second: 36
	// Only once this second: 37
}
