package backoff

import (
	"context"
	"math"
	"time"
)

const (
	defaultInterval     = 500 * time.Millisecond
	defaultJitterFactor = 0.5
	defaultThreshold    = 15 * time.Minute
)

func NewExponential(options ...Option) *Exponential {
	interval := defaultInterval
	jitterFactor := defaultJitterFactor
	maxRetries := -1
	threshold := defaultThreshold
	factor := float64(2)
	for _, o := range options {
		switch o.Name() {
		case optkeyInterval:
			interval = o.Value().(time.Duration)
		case optkeyJitterFactor:
			jitterFactor = o.Value().(float64)
		case optkeyMaxRetries:
			maxRetries = o.Value().(int)
		case optkeyThreshold:
			threshold = o.Value().(time.Duration)
		}
	}

	return &Exponential{
		factor:       factor,
		interval:     interval,
		jitterFactor: jitterFactor,
		maxRetries:   maxRetries,
		threshold:    threshold,
	}
}

func (p *Exponential) Start(ctx context.Context) (Backoff, CancelFunc) {
	b := &exponentialBackoff{
		baseBackoff:  newBaseBackoff(ctx, p.maxRetries),
		policy: p,
	}

	return b, CancelFunc(b.cancelLocked)
}

func (b *exponentialBackoff) Next() <-chan struct{} {
	d := b.delayForAttempt(b.attempt)
	b.attempt++
	time.AfterFunc(d, b.fire)
	return b.next
}

func (b *exponentialBackoff) delayForAttempt(attempt float64) time.Duration {
	minf := float64(b.policy.interval)
	durf := minf * math.Pow(b.policy.factor, attempt)

	dur := time.Duration(durf)

	return dur
}
