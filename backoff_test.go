package backoff_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	backoff "github.com/lestrrat/go-backoff"
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
