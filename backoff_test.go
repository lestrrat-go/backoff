package backoff_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/lestrrat-go/backoff/v2"
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

func isInErrorRange(expected, observed, margin time.Duration) bool {
	return expected+margin > observed &&
		observed > expected-margin
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
		ig := backoff.NewExponentialInterval(
			backoff.WithMaxInterval(time.Second),
			backoff.WithJitterFactor(0.02),
		)

		testcases := []struct {
			Base time.Duration
		}{
			{Base: 500 * time.Millisecond},
			{Base: 750 * time.Millisecond},
			{Base: time.Second},
		}

		for i := 0; i < 10; i++ {
			dur := ig.Next()
			var base time.Duration
			if i > 2 {
				base = testcases[2].Base
			} else {
				base = testcases[i].Base
			}

			min := int64(float64(base) * 0.98)
			max := int64(float64(base) * 1.05) // should be 1.02, but give it a bit of leeway
			t.Logf("max = %s, min = %s", time.Duration(max), time.Duration(min))
			if !assert.GreaterOrEqual(t, int64(dur), min, "value should be greater than minimum") {
				return
			}
			if !assert.GreaterOrEqual(t, max, int64(dur), "value should be less than maximum") {
				return
			}

		}
	})
	t.Run("Back off, no jitter", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		// These values are truncated to milliseconds, to make comparisons easier
		expected := []float64{
			0, 0.5, 0.7, 1.1, 1.6, 2.5, 3.7,
		}
		p := backoff.Exponential()
		count := 0
		prev := time.Now()
		b := p.Start(ctx)
		for backoff.Continue(b) {
			now := time.Now()
			d := now.Sub(prev)
			d = d - d%(100*time.Millisecond)

			// Allow a flux of 100ms
			expectedDuration := time.Duration(expected[count] * float64(time.Second))
			if !assert.True(t, isInErrorRange(expectedDuration, d, 100*time.Millisecond), `observed duration (%s) should be whthin error range (expected = %s, range = %s)`, d, expectedDuration, 100*time.Millisecond) {
				return
			}
			count++
			if count == len(expected)-1 {
				break
			}
			prev = now
		}
	})
}

func TestConcurrent(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	t.Parallel()

	// Does not test anything useful, just puts it under stress
	testcases := []struct {
		Policy backoff.Policy
		Name   string
	}{
		{Name: "Null", Policy: backoff.Null()},
		{Name: "Exponential", Policy: backoff.Exponential(backoff.WithMultiplier(0.01), backoff.WithMinInterval(time.Millisecond))},
	}

	const max = 50
	for _, tc := range testcases {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			var wg sync.WaitGroup
			wg.Add(max)
			for i := 0; i < max; i++ {
				go func(wg *sync.WaitGroup, b backoff.Policy) {
					defer wg.Done()
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					c := b.Start(ctx)
					for backoff.Continue(c) {
						fmt.Fprintf(ioutil.Discard, `Writing to the ether...`)
					}
				}(&wg, tc.Policy)
			}
			wg.Wait()
		})
	}
}

func TestConstantWithJitter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	p := backoff.Constant(
		backoff.WithInterval(300*time.Millisecond),
		backoff.WithJitterFactor(0.50),
		backoff.WithMaxRetries(999),
	)
	c := p.Start(ctx)

	prev := time.Now()
	var retries int
	for backoff.Continue(c) {
		t.Logf("%s backoff.Continue", time.Now())

		// make sure that we've executed this in more or less 300ms ± 50%
		retries++
		if retries > 1 {
			d := time.Since(prev)

			// if the duration becomes out of the range values by jitter, it breaks loop
			if (150*time.Millisecond <= d && d < 250*time.Millisecond) ||
				(350*time.Millisecond < d && d <= 450*time.Millisecond) {
				break
			}
		}
		prev = time.Now()
	}

	// initial + 999 retries = 1000
	if !assert.NotEqual(t, 1000, retries, `should not have retried 1000 times; if the # of retries reaches 1000, probably jitter doesn't work'`) {
		return
	}
}
