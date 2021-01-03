package backoff_test

import (
	"context"
	"testing"
	"time"

	"github.com/lestrrat-go/backoff"
	"github.com/stretchr/testify/assert"
)

func TestNull(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	p := backoff.Null()
	c := p.Start(ctx)

	var retries int
	for backoff.Continue(c) {
		t.Logf("%s backoff.Continue", time.Now())
		retries++
	}
	if !assert.Equal(t, 1, retries, `should have done 1 retries`) {
		return
	}
}

func TestConstant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	p := backoff.Constant(
		backoff.WithInterval(300*time.Millisecond),
		backoff.WithMaxRetries(4),
	)
	c := p.Start(ctx)

	prev := time.Now()
	var retries int
	for backoff.Continue(c) {
		t.Logf("%s backoff.Continue", time.Now())

		// make sure that we've executed this in more or less 300ms
		retries++
		if retries > 1 {
			d := time.Since(prev)
			if !assert.True(t, 350*time.Millisecond >= d && d >= 250*time.Millisecond, `timing is about 300ms (%s)`, d) {
				return
			}
		}
		prev = time.Now()
	}

	// initial + 4 retries = 5
	if !assert.Equal(t, 5, retries, `should have retried 5 times`) {
		return
	}
}

func TestExponential(t *testing.T) {
	t.Run("Interval generator", func(t *testing.T) {
		expected := []float64{
			0.5, 0.75, 1.125, 1.6875, 2.53125, 3.796875,
		}
		ig := backoff.NewExponentialInterval()
		for i := 0; i < len(expected); i++ {
			if !assert.Equal(t, time.Duration(float64(time.Second)*expected[i]), ig.Next(), `interval for iteration %d`, i) {
				return
			}
		}
	})
	t.Run("Jitter", func(t *testing.T) {
		ig := backoff.NewExponentialInterval(backoff.WithJitterFactor(0.02))
		for i := 0; i < 10; i++ {
			t.Logf("%s", ig.Next())
		}
	})
}
