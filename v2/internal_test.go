package backoff

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptionPassing(t *testing.T) {
	cOptions := []ControllerOption{
		WithMaxRetries(9999999999999),
	}
	igOptions := []ExponentialOption{
		WithJitterFactor(0.99),
		WithMaxInterval(24 * time.Hour),
		WithMinInterval(time.Nanosecond),
		WithMultiplier(99999),
		WithRNG(rand.New(rand.NewSource(time.Now().UnixNano()))),
	}

	merged := igOptions
	for _, option := range cOptions {
		merged = append(merged, option.(ExponentialOption))
	}
	p := NewExponentialPolicy(merged...)

	if !assert.Equal(t, cOptions, p.cOptions) {
		return
	}
	if !assert.Equal(t, igOptions, p.igOptions) {
		return
	}
}
