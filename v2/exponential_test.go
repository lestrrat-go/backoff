package backoff

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExponentialIntervalWithDefaultOptions(t *testing.T) {
	p := NewExponentialInterval()

	assert.Equal(t, 0.0, p.jitterFactor)
	assert.Equal(t, defaultMaxInterval, p.maxInterval)
	assert.Equal(t, defaultMinInterval, p.minInterval)
	assert.Equal(t, defaultMultiplier, p.multiplier)
	assert.Nil(t, p.rng)
}

func TestNewExponentialIntervalWithCustomOptions(t *testing.T) {
	jitter := 0.99
	maxInterval := 24 * time.Hour
	minInterval := time.Nanosecond
	multiplier := float64(99999)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := NewExponentialInterval(
		WithJitterFactor(jitter),
		WithMaxInterval(maxInterval),
		WithMinInterval(minInterval),
		WithMultiplier(multiplier),
		WithRNG(rng),
	)

	assert.Equal(t, jitter, p.jitterFactor)
	assert.Equal(t, maxInterval, time.Duration(p.maxInterval))
	assert.Equal(t, minInterval, time.Duration(p.minInterval))
	assert.Equal(t, multiplier, p.multiplier)
	assert.Equal(t, rng, p.rng)
}

func TestNewExponentialIntervalWithOnlyJitterOptions(t *testing.T) {
	jitter := 0.99
	p := NewExponentialInterval(
		WithJitterFactor(jitter),
	)

	assert.Equal(t, jitter, p.jitterFactor)
	assert.NotNil(t, p.rng, "should be generated automatically")
}
