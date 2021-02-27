package backoff_test

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/lestrrat-go/backoff/v2"
)

func TestLeak(t *testing.T) {
	beforeGoroutine := runtime.NumGoroutine()
	var wg sync.WaitGroup
	tasks := 100
	wg.Add(tasks)
	for range make([]struct{}, tasks) {
		go func() {
			defer wg.Done()
			exp := backoff.Exponential()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := exp.Start(ctx)
			for backoff.Continue(c) {
				time.Sleep(1300 * time.Millisecond)
				cancel()
			}
		}()
	}
	wg.Wait()
	afterGoroutine := runtime.NumGoroutine()
	if afterGoroutine > beforeGoroutine+10 {
		t.Errorf("goroutines seem to be leaked. before: %d, after: %d", beforeGoroutine, afterGoroutine)
	}
}
