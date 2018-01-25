package backoff_test

import (
	"context"
	"errors"
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
		if !assert.True(t, elapsed < 600 * time.Millisecond, `total elapsed time should be less than 600 ms (%s)`, elapsed) {
			return
		}

		if !assert.Equal(t, 10, count, `fn should have been called 10 times`) {
			return
		}
	})
}
