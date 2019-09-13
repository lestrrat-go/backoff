package backoff_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	backoff "github.com/lestrrat-go/backoff"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	count := 0
	fn := func(ctx context.Context) error {
		if count++; count%10 == 0 {
			return nil
		}
		return errors.New(`dummy`)
	}

	t.Run("succeed on 10th try", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		p := backoff.NewExponential(
			backoff.WithInterval(time.Millisecond),
			backoff.WithJitterFactor(0),
			backoff.WithFactor(2),
		)

		start := time.Now()
		err := backoff.Retry(ctx, p, backoff.ExecuteFunc(fn))
		if !assert.NoError(t, err, `backoff.Retry should succeed`) {
			return
		}

		// 1 + 2 + 4 + 8 ... + 256 == 511 ms. with everything going on, we should
		// never exceed  600 ms
		elapsed := time.Since(start)
		if !assert.True(t, elapsed < 600*time.Millisecond, `total elapsed time should be less than 600 ms (%s)`, elapsed) {
			return
		}

		if !assert.Equal(t, 10, count, `fn should have been called 10 times`) {
			return
		}
	})
}

func TestRetryExponentialParallel(t *testing.T) {
	p := backoff.NewExponential(
		backoff.WithInterval(time.Millisecond),
		backoff.WithJitterFactor(0),
		backoff.WithFactor(2),
	)

	// Timeout here needs to be handled carefully, as parallel t.Run()
	// will fallthrough out of this function, and any defered cancel()
	// statements will be called before the tests have a chance to run
	// properly. So instead of using WithTimeout() then defer the cancel
	// function like we normally do, we wait calling cancel() until
	// all of the subtests are done
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	var wg sync.WaitGroup

	// Just run a bunch of goroutines, check for races
	for i := 0; i < 10; i++ {
		count := 0
		fn := backoff.ExecuteFunc(func(ctx context.Context) error {
			if count++; count%10 == 0 {
				return nil
			}
			return errors.New(`dummy`)
		})

		wg.Add(1)
		t.Run(fmt.Sprintf("goroutine %02d", i), func(t *testing.T) {
			defer wg.Done()
			t.Parallel()
			err := backoff.Retry(ctx, p, fn)
			if !assert.NoError(t, err, `backoff.Retry should succeed`) {
				return
			}
		})
	}

	go func() {
		wg.Wait()
		cancel() // okay, now we can cancel
	}()
}

func TestGHIssue1(t *testing.T) {
	makeOptions := func(options ...backoff.Option) []backoff.Option {
		var newOptions []backoff.Option

		newOptions = append(newOptions, backoff.WithInterval(10*time.Millisecond))
		newOptions = append(newOptions, backoff.WithMaxRetries(1))
		newOptions = append(newOptions, options...)
		return newOptions
	}

	run := func(ctx context.Context, options ...backoff.Option) {
		policy := backoff.NewExponential(options...)
		b, cancel := policy.Start(ctx)
		defer cancel()

		for {
			select {
			case <-b.Done():
				return
			case <-b.Next():
			}
		}
	}

	tests := []struct {
		options []backoff.Option
	}{
		{
			options: makeOptions(),
		},
		{
			options: makeOptions(backoff.WithMaxRetries(20)),
		},
	}

	for _, test := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		run(ctx, test.options...)
		cancel()
	}
}

func TestMaxElapsedTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	policy := backoff.NewConstant(100*time.Millisecond, backoff.WithMaxElapsedTime(time.Second))
	b, bcancel := policy.Start(ctx)
	defer bcancel()

	var count int
LOOP:
	for {
		select {
		case <-b.Next():
			count++
		case <-b.Done():
			break LOOP
		case <-ctx.Done():
			t.Errorf("context expired before backoff")
			return
		}
	}
	if !assert.True(t, count > 5, "we should have at least a few iterations") {
		return
	}
}

func TestGHIssue6(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t.Run("constant backoff", func(t *testing.T) {
		backoffPolicy := backoff.NewConstant(10 * time.Millisecond)
		b, cancel := backoffPolicy.Start(ctx)
		defer cancel()

		seen := make(map[time.Time]struct{})
		for backoff.Continue(b) {
			// Record execution in millisecond granularity for testing
			now := time.Now().Truncate(time.Millisecond)
			if _, ok := seen[now]; !assert.False(t, ok, `should not fire in the same millisecond`) {
				return
			}
			seen[now] = struct{}{}
		}
	})
	t.Run("exponential backoff", func(t *testing.T) {
		backoffPolicy := backoff.NewExponential(
			backoff.WithInterval(10*time.Millisecond),
			backoff.WithJitterFactor(0),
			backoff.WithFactor(2),
		)

		b, cancel := backoffPolicy.Start(ctx)
		defer cancel()

		seen := make(map[time.Time]struct{})
		for backoff.Continue(b) {
			// Record execution in millisecond granularity for testing
			now := time.Now().Truncate(time.Millisecond)
			if _, ok := seen[now]; !assert.False(t, ok, `should not fire in the same millisecond`) {
				return
			}
			seen[now] = struct{}{}
		}
	})
}

func TestRetryForever(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	count := 0
	go func() {
		defer close(done)
		policy := backoff.NewExponential(
			backoff.WithRetryForever(),
			backoff.WithInterval(time.Millisecond),
		)
		b, cancel := policy.Start(ctx)
		defer cancel()

		tick := time.NewTicker(time.Millisecond)
		for backoff.Continue(b) {
			count++
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
			}
		}

		// if we got here, it means that we did not retry forever
		t.Errorf("Bailed out of backoff.Continue(b)")
	}()

	timer := time.NewTimer(time.Second)
	<-timer.C

	cancel()

	// wait for max 1 second for the goroutine to come back
	timeout := time.NewTimer(time.Second)
	select {
	case <-done:
		if assert.True(t, count >= 10, "we should have executed more than 10 times, but executed %d", count) {
			return
		}
	case <-timeout.C:
		t.Error(`goroutine did not come back in time`)
	}
}
