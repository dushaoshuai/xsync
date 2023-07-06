# About

[![Go Reference](https://pkg.go.dev/badge/github.com/dushaoshuai/xsync.svg)](https://pkg.go.dev/github.com/dushaoshuai/xsync)

Some Go concurrency utilities.

# Download/Install

```shell
go get -u github.com/dushaoshuai/xsync
```

# Usage

## Onces

Onces is an object that will try to (if asked) perform a successful action only
if a specified interval has elapsed since the last successful action.

Example:

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/dushaoshuai/xsync"
)

func main() {
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
```
