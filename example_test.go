package xsync_test

import (
	"fmt"
	"sync"
	"time"

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
	// Only once this second: 44
	// Only once this second: 45
}
