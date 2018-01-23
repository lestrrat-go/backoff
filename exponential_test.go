package backoff

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponential(t *testing.T) {
	p := NewExponential(WithInterval(time.Second))
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
