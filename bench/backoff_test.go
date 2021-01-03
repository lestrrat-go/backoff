// +build bench

package backoff_test

import (
	"context"
	"errors"
	"testing"

	cenkalti "github.com/cenkalti/backoff"
	lestrrat "github.com/lestrrat-go/backoff"
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

	var sink int
	b.Run("cenkalti", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			backoff := cenkalti.NewExponentialBackOff()
			b.StartTimer()
			cenkalti.Retry(func() error {
				sink = fn()
				return errors.New(`dummy`)
			}, cenkalti.WithMaxRetries(backoff, 5))
		}
	})
	b.Run("lestrrat", func(b *testing.B) {
		b.StopTimer()
		policy := lestrrat.NewExponential(lestrrat.WithMaxRetries(5), lestrrat.WithFactor(1.2))
		for i := 0; i < b.N; i++ {
			b.StartTimer()
			backoff, cancel := policy.Start(context.Background())
		MAIN:
			for {
				fn()
				select {
				case <-backoff.Done():
					break MAIN
				case <-backoff.Next():
				}
			}
			cancel()
			b.StopTimer()
		}
	})
}
