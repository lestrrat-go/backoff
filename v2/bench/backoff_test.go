package backoff_test

import (
	"context"
	"errors"
	"testing"

	cenkalti "github.com/cenkalti/backoff"
	lestrrat "github.com/lestrrat-go/backoff/v2"
)

func Benchmark(b *testing.B) {
	// This is a dummy function
	fn := func() int {
		var v int
		for i := 1; i <= 10; i++ {
			v += i
		}
		return v
	}

	b.Run("cenkalti", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			backoff := cenkalti.NewExponentialBackOff()
			b.StartTimer()
			cenkalti.Retry(func() error {
				_ = fn()
				return errors.New(`dummy`)
			}, cenkalti.WithMaxRetries(backoff, 5))
		}
	})
	b.Run("lestrrat", func(b *testing.B) {
		b.StopTimer()
		policy := lestrrat.Exponential(lestrrat.WithMaxRetries(5), lestrrat.WithJitterFactor(1.2))
		for i := 0; i < b.N; i++ {
			b.StartTimer()
			backoff := policy.Start(context.Background())
		MAIN:
			for {
				fn()
				select {
				case <-backoff.Done():
					break MAIN
				case <-backoff.Next():
					_ = fn()
				}
			}
			b.StopTimer()
		}
	})
}
