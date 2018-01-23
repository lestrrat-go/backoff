package backoff

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponential(t *testing.T) {
	p := NewExponential(WithInterval(time.Second), WithJitterFactor(0))
	b, cancel := p.Start(context.Background())
	defer cancel()

	const attempts = 10
	var durs = []time.Duration{
		time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		16 * time.Second,
		32 * time.Second,
		64 * time.Second,
		128 * time.Second,
		256 * time.Second,
		512 * time.Second,
	}
	for i := 0; i < 10; i++ {
		dur := b.(*exponentialBackoff).delayForAttempt(float64(i))
		if !assert.Equal(t, dur, durs[i], `delays should match`) {
			return
		}
	}
}

func TestExponentialWithJitter(t *testing.T) {
	const jitter = 0.2
	p := NewExponential(WithInterval(time.Second), WithJitterFactor(jitter))
	b, cancel := p.Start(context.Background())
	defer cancel()

	const attempts = 10
	var durs = []time.Duration{
		time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		16 * time.Second,
		32 * time.Second,
		64 * time.Second,
		128 * time.Second,
		256 * time.Second,
		512 * time.Second,
	}
	for i := 0; i < 10; i++ {
		dur := b.(*exponentialBackoff).delayForAttempt(float64(i))

		durf := float64(dur)
		expectedf := float64(durs[i])
		delta := expectedf * jitter
		max := expectedf + delta
		min := expectedf - delta
		if !assert.True(t, min <= durf && max >= durf, `delays should be between %f and %f, got %f`, min, max, durf) {
			return
		}
	}
}
