package backoff_test

import (
	"context"
	"runtime"
	"sync"
	"testing"

	"github.com/lestrrat-go/backoff/v2"
)

func TestLeakNull(t *testing.T) {
	beforeGoroutine := runtime.NumGoroutine()
	var wg sync.WaitGroup
	tasks := 100
	wg.Add(tasks)
	for range make([]struct{}, tasks) {
		go func() {
			defer wg.Done()
			null := backoff.Null()
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			c := null.Start(ctx)
			for backoff.Continue(c) {
				return
			}
		}()
	}
	wg.Wait()
	afterGoroutine := runtime.NumGoroutine()
	if afterGoroutine > beforeGoroutine+10 {
		t.Errorf("goroutines seem to be leaked. before: %d, after: %d", beforeGoroutine, afterGoroutine)
	}
}
